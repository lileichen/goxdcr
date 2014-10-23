package parts

import (
	"encoding/binary"
	"errors"
	"fmt"
	common "github.com/Xiaomei-Zhang/couchbase_goxdcr/common"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr/log"
	part "github.com/Xiaomei-Zhang/couchbase_goxdcr/part"
	base "github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/base"
	gen_server "github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/gen_server"
	"github.com/Xiaomei-Zhang/couchbase_goxdcr_impl/utils"
	mc "github.com/couchbase/gomemcached"
	mcc "github.com/couchbase/gomemcached/client"
	"io"
	"math"
	"math/rand"
	"net"
	"reflect"
	"sync"
	"time"
)

type XMEM_MODE int

const (
	Batch_XMEM        XMEM_MODE = iota
	Asynchronous_XMEM XMEM_MODE = iota
)

//configuration settings for XmemNozzle
const (
	//configuration param names
	XMEM_SETTING_BATCHCOUNT            = "batch_count"
	XMEM_SETTING_BATCHSIZE             = "batch_size"
	XMEM_SETTING_MODE                  = "mode"
	XMEM_SETTING_NUMOFRETRY            = "max_retry"
	XMEM_SETTING_TIMEOUT               = "timeout"
	XMEM_SETTING_BATCH_EXPIRATION_TIME = "batch_expiration_time"

	//default configuration
	default_batchcount          int           = 500
	default_batchsize           int           = 2048
	default_mode                XMEM_MODE     = Asynchronous_XMEM
	default_numofretry          int           = 6
	default_timeout             time.Duration = 100 * time.Millisecond
	default_dataChannelSize                   = 5000
	default_batchExpirationTime               = 100 * time.Millisecond
)

const (
	SET_WITH_META    = mc.CommandCode(0xa2)
	DELETE_WITH_META = mc.CommandCode(0xa8)
)

var xmem_setting_defs base.SettingDefinitions = base.SettingDefinitions{XMEM_SETTING_BATCHCOUNT: base.NewSettingDef(reflect.TypeOf((*int)(nil)), false),
	XMEM_SETTING_BATCHSIZE:             base.NewSettingDef(reflect.TypeOf((*int)(nil)), false),
	XMEM_SETTING_MODE:                  base.NewSettingDef(reflect.TypeOf((*XMEM_MODE)(nil)), false),
	XMEM_SETTING_NUMOFRETRY:            base.NewSettingDef(reflect.TypeOf((*int)(nil)), false),
	XMEM_SETTING_TIMEOUT:               base.NewSettingDef(reflect.TypeOf((*time.Duration)(nil)), false),
	XMEM_SETTING_BATCH_EXPIRATION_TIME: base.NewSettingDef(reflect.TypeOf((*time.Duration)(nil)), false)}

/************************************
/* struct bufferedMCRequest
*************************************/

type bufferedMCRequest struct {
	req          *mc.MCRequest
	sent_time    time.Time
	num_of_retry int
	err          error
	reservation  int
}

func newBufferedMCRequest(request *mc.MCRequest, reservationNum int) *bufferedMCRequest {
	return &bufferedMCRequest{req: request,
		sent_time:    time.Now(),
		num_of_retry: 0,
		err:          nil,
		reservation:  reservationNum}
}

/***********************************************************
/* struct requestBuffer
/* This is used to buffer the sent but yet confirmed data
************************************************************/
type requestBuffer struct {
	slots           []*bufferedMCRequest /*slots to store the data*/
	sequences       []uint16
	empty_slots_pos chan uint16 /*empty slot pos in the buffer*/
	size            uint16      /*the size of the buffer*/
	notifych        chan bool   /*notify channel is set when the buffer is empty*/
	logger          *log.CommonLogger
}

func newReqBuffer(size uint16, notifychan chan bool, logger *log.CommonLogger) *requestBuffer {
	logger.Debugf("Create a new request buffer of size %d\n", size)
	buf := &requestBuffer{
		make([]*bufferedMCRequest, size, size),
		make([]uint16, size),
		make(chan uint16, size),
		size,
		notifychan,
		logger}

	logger.Debug("Slots is initialized")

	//initialize the empty_slots_pos
	buf.initializeEmptySlotPos()

	logger.Debugf("new request buffer of size %d is created\n", size)
	return buf
}

func (buf *requestBuffer) initializeEmptySlotPos() error {
	for i := 0; i < int(buf.size); i++ {
		buf.empty_slots_pos <- uint16(i)
		buf.sequences[i] = 0
	}

	return nil
}

func (buf *requestBuffer) validatePos(pos uint16) (err error) {
	err = nil
	if pos < 0 || int(pos) >= len(buf.slots) {
		buf.logger.Error("Invalid slot index")
		err = errors.New("Invalid slot index")
	}
	return
}

//slot allow caller to get hold of the content in the slot without locking the slot
//@pos - the position of the slot
func (buf *requestBuffer) slot(pos uint16) (*mc.MCRequest, error) {
	buf.logger.Debugf("Getting the content in slot %d\n", pos)

	err := buf.validatePos(pos)
	if err != nil {
		return nil, err
	}

	req := buf.slots[pos]

	if req == nil {
		return nil, nil
	} else {
		return req.req, nil
	}

}

//modSlot allow caller to do book-keeping on the slot, like updating num_of_retry, err
//@pos - the position of the slot
//@modFunc - the callback function which is going to update the slot
func (buf *requestBuffer) modSlot(pos uint16, modFunc func(req *bufferedMCRequest, p uint16) bool) (bool, error) {
	var err error = nil
	err = buf.validatePos(pos)
	if err != nil {
		return false, err
	}

	req := buf.slots[pos]

	var modified bool

	if req != nil && req.req != nil {
		modified = modFunc(req, pos)
	} else {
		modified = false
	}
	return modified, err
}

//evictSlot allow caller to empty the slot
//@pos - the position of the slot
func (buf *requestBuffer) evictSlot(pos uint16) error {

	err := buf.validatePos(pos)
	if err != nil {
		return err
	}
	req := buf.slots[pos]
	buf.slots[pos] = nil

	if req != nil {
		buf.empty_slots_pos <- pos

		//increase sequence
		if buf.sequences[pos]+1 > 65535 {
			buf.sequences[pos] = 0
		} else {
			buf.sequences[pos] = buf.sequences[pos] + 1
		}

		buf.logger.Debugf("buf.seqence[%v]=%v\n", pos, buf.sequences[pos])
		buf.logger.Debugf("wait_for_resp=%vs\n", time.Since(req.sent_time).Seconds())

		if len(buf.empty_slots_pos) == int(buf.size) {
			if buf.notifych != nil {
				buf.notifych <- true
				buf.logger.Debug("buffer is empty, notify")
			} else {
				buf.logger.Debug("buffer is empty, no notify channel is specified though")
			}
		}
	}
	return nil

}

//availableSlotIndex returns a position number of an empty slot
func (buf *requestBuffer) reserveSlot() (error, uint16, int) {
	buf.logger.Debugf("slots chan length=%d\n", len(buf.empty_slots_pos))
	index := <-buf.empty_slots_pos

	var reservation_num int

	//non blocking
	//generate a random number
	reservation_num = rand.Int()
	req := newBufferedMCRequest(nil, reservation_num)
	buf.slots[index] = req
	return nil, uint16(index), reservation_num
}

func (buf *requestBuffer) cancelReservation(index uint16, reservation_num int) error {

	err := buf.validatePos(index)
	if err == nil {
		buf.empty_slots_pos <- index

		req := buf.slots[index]
		var reservation_num int

		//non blocking

		if req.reservation == reservation_num {
			req = nil
		} else {
			err = errors.New("Cancel reservation failed, reservation number doesn't match")
		}
	}
	return err
}

func (buf *requestBuffer) enSlot(pos uint16, req *mc.MCRequest, reservationNum int) error {
	buf.logger.Debugf("enSlot: pos=%d\n", pos)

	err := buf.validatePos(pos)
	if err != nil {
		return err
	}
	r := buf.slots[pos]

	if r == nil {
		buf.slots[pos] = newBufferedMCRequest(nil, 0)
	} else {
		if r.reservation != reservationNum {
			buf.logger.Errorf("Can't enSlot %d, doesn't have the reservation, %v", pos, r)
			return errors.New(fmt.Sprintf("Can't enSlot %d, doesn't have the reservation", pos))
		}
		r.req = req
	}
	buf.logger.Debugf("slot %d is occupied\n", pos)
	return nil
}

func (buf *requestBuffer) bufferSize() uint16 {
	return buf.size
}

/************************************
/* struct xmemConfig
*************************************/
type xmemConfig struct {
	maxCount int
	maxSize  int
	//the duration to wait for the batch-sending to finish
	batchtimeout        time.Duration
	batchExpirationTime time.Duration
	maxRetry            int
	mode                XMEM_MODE
	connectStr          string
	bucketName          string
	password            string
	logger              *log.CommonLogger
}

func newConfig(logger *log.CommonLogger) xmemConfig {
	return xmemConfig{maxCount: default_batchcount,
		maxSize:             default_batchsize,
		batchtimeout:        default_timeout,
		batchExpirationTime: default_batchExpirationTime,
		maxRetry:            default_numofretry,
		mode:                default_mode,
		connectStr:          "",
		bucketName:          "",
		password:            "",
	}

}

func (config *xmemConfig) initializeConfig(settings map[string]interface{}) error {
	err := utils.ValidateSettings(xmem_setting_defs, settings, config.logger)
	if val, ok := settings[XMEM_SETTING_BATCHSIZE]; ok {
		config.maxSize = val.(int)
	}
	if val, ok := settings[XMEM_SETTING_BATCHCOUNT]; ok {
		config.maxCount = val.(int)
	}
	if val, ok := settings[XMEM_SETTING_TIMEOUT]; ok {
		config.batchtimeout = val.(time.Duration)
	}
	if val, ok := settings[XMEM_SETTING_NUMOFRETRY]; ok {
		config.maxRetry = val.(int)
	}
	if val, ok := settings[XMEM_SETTING_MODE]; ok {
		config.mode = val.(XMEM_MODE)
	}
	if val, ok := settings[XMEM_SETTING_BATCH_EXPIRATION_TIME]; ok {
		config.batchExpirationTime = val.(time.Duration)
	}
	return err
}

/************************************
/* struct xmemBatch
*************************************/
type xmemBatch struct {
	curCount       int
	curSize        int
	capacity_count int
	capacity_size  int
	start_time     time.Time
	frozen         bool
	logger         *log.CommonLogger
}

func newXmemBatch(cap_count int, cap_size int, logger *log.CommonLogger) *xmemBatch {
	return &xmemBatch{
		curCount:       0,
		curSize:        0,
		capacity_count: cap_count,
		capacity_size:  cap_size,
		logger:         logger}
}

func (b *xmemBatch) accumuBatch(size int) bool {
	var ret bool = true

	b.curCount++
	if b.curCount == 1 {
		b.start_time = time.Now()
	}

	b.curSize += size
	if b.curCount < b.capacity_count && b.curSize < b.capacity_size*1000 {
		ret = false
	}
	return ret
}

func (b *xmemBatch) count() int {
	return b.curCount
}

func (b *xmemBatch) size() int {
	return b.curSize

}

/************************************
/* struct XmemNozzle
*************************************/
type XmemNozzle struct {

	//parent inheritance
	gen_server.GenServer
	part.AbstractPart

	bOpen      bool
	lock_bOpen sync.RWMutex

	//data channel to accept the incoming data
	dataChan chan *mc.MCRequest

	//memcached client connected to the target bucket
	memClient *mcc.Client

	//configurable parameter
	config xmemConfig

	//queue for ready batches
	batches_ready chan *xmemBatch

	//batch to be accumulated
	batch *xmemBatch

	//channel to signal if batch can be transitioned to batches_ready
	batch_move_ch chan bool

	//control channel
	sendNowCh chan bool

	childrenWaitGrp sync.WaitGroup

	//buffer for the sent, but not yet confirmed data
	buf *requestBuffer

	sender_finch   chan bool
	receiver_finch chan bool
	checker_finch  chan bool
	send_allow_ch  chan bool

	counter_sent     int
	counter_received int
	start_time       time.Time
}

func NewXmemNozzle(id string,
	connectString string,
	bucketName string,
	password string,
	logger_context *log.LoggerContext) *XmemNozzle {

	//callback functions from GenServer
	var msg_callback_func gen_server.Msg_Callback_Func
	//	var behavior_callback_func gen_server.Behavior_Callback_Func
	var exit_callback_func gen_server.Exit_Callback_Func
	var error_handler_func gen_server.Error_Handler_Func

	var isStarted_callback_func part.IsStarted_Callback_Func

	server := gen_server.NewGenServer(&msg_callback_func,
		nil, &exit_callback_func, &error_handler_func, logger_context, "XmemNozzle")
	isStarted_callback_func = server.IsStarted
	part := part.NewAbstractPartWithLogger(id, &isStarted_callback_func, server.Logger())

	xmem := &XmemNozzle{GenServer: server, /*gen_server.GenServer*/
		AbstractPart:     part,                       /*part.AbstractPart*/
		bOpen:            true,                       /*bOpen	bool*/
		lock_bOpen:       sync.RWMutex{},             /*lock_bOpen	sync.RWMutex*/
		dataChan:         nil,                        /*dataChan*/
		memClient:        nil,                        /*memClient*/
		config:           newConfig(server.Logger()), /*config	xmemConfig*/
		batches_ready:    make(chan *xmemBatch, 100), /*batches_ready chan *xmemBatch*/
		batch:            nil,                        /*batch		  *xmemBatch*/
		batch_move_ch:    nil,
		sendNowCh:        make(chan bool, 1), /*sendNowCh chan bool*/
		childrenWaitGrp:  sync.WaitGroup{},   /*childrenWaitGrp sync.WaitGroup*/
		buf:              nil,                /*buf	requestBuffer*/
		receiver_finch:   make(chan bool),    /*receiver_finch chan bool*/
		checker_finch:    make(chan bool),    /*checker_finch chan bool*/
		sender_finch:     make(chan bool),
		send_allow_ch:    make(chan bool, 1), /*send_allow_ch chan bool*/
		counter_sent:     0,
		counter_received: 0}

	xmem.config.connectStr = connectString
	xmem.config.bucketName = bucketName
	xmem.config.password = password

	msg_callback_func = nil
	exit_callback_func = xmem.onExit
	error_handler_func = xmem.handleGeneralError
	return xmem

}

func (xmem *XmemNozzle) Open() error {
	if !xmem.bOpen {
		xmem.bOpen = true

	}
	return nil
}

func (xmem *XmemNozzle) Close() error {
	if xmem.bOpen {
		xmem.bOpen = false
	}
	return nil
}

func (xmem *XmemNozzle) Start(settings map[string]interface{}) error {
	xmem.Logger().Info("Xmem starting ....")
	err := xmem.initialize(settings)
	xmem.Logger().Info("....Finish initializing....")
	if err == nil {
		xmem.childrenWaitGrp.Add(1)
		go xmem.receiveResponse(xmem.memClient, xmem.receiver_finch, &xmem.childrenWaitGrp)

		xmem.childrenWaitGrp.Add(1)
		go xmem.check(xmem.checker_finch, &xmem.childrenWaitGrp)

		xmem.childrenWaitGrp.Add(1)

		if xmem.config.mode == Batch_XMEM {
			go xmem.processData_batch(xmem.sender_finch, &xmem.childrenWaitGrp)
		} else {
			go xmem.processData_async(xmem.sender_finch, &xmem.childrenWaitGrp)
		}
	}
	xmem.start_time = time.Now()
	if err == nil {
		err = xmem.Start_server()
	}

	xmem.Logger().Info("Xmem nozzle is started")
	return err
}

func (xmem *XmemNozzle) getReadyToShutdown() {
	xmem.Logger().Debug("Waiting for data is drained")

	for {
		if len(xmem.dataChan) == 0 && len(xmem.batches_ready) == 0 {
			xmem.Logger().Debug("Ready to stop")
			break
		} else if len(xmem.batches_ready) == 0 && xmem.batch.count() > 0 {
			xmem.batchReady()

		} else {
			xmem.Logger().Debugf("%d in data channel, %d batches ready, % data in current batch \n", len(xmem.dataChan), len(xmem.batches_ready), xmem.batch.count())
		}
	}

	close(xmem.batches_ready)
}

func (xmem *XmemNozzle) Stop() error {
	xmem.Logger().Infof("Stop XmemNozzle %v\n", xmem.Id())
	xmem.getReadyToShutdown()

	conn := xmem.memClient.Hijack()
	xmem.memClient = nil
	conn.(*net.TCPConn).SetReadDeadline(time.Now())

	xmem.Logger().Debugf("XmemNozzle %v processed %v items\n", xmem.Id(), xmem.counter_sent)
	err := xmem.Stop_server()

	conn.(*net.TCPConn).SetReadDeadline(time.Date(1, time.January, 0, 0, 0, 0, 0, time.UTC))
	xmem.Logger().Debugf("XmemNozzle %v is stopped\n", xmem.Id())
	return err
}

func (xmem *XmemNozzle) IsOpen() bool {
	ret := xmem.bOpen
	return ret
}

func (xmem *XmemNozzle) batchReady() {
	//move the batch to ready batches channel
	if xmem.batch.count() > 0 {
		<-xmem.batch_move_ch
		defer func() {
			xmem.batch_move_ch <- true
			xmem.Logger().Infof("End moving batch, %v batches ready\n", len(xmem.batches_ready))
		}()
		xmem.Logger().Infof("move the batch (count=%d) ready queue\n", xmem.batch.count())
		select {
		case xmem.batches_ready <- xmem.batch:

			xmem.Logger().Debugf("There are %d batches in ready queue\n", len(xmem.batches_ready))
			xmem.initNewBatch()
		default:
		}
	}

}

func (xmem *XmemNozzle) Receive(data interface{}) error {
	xmem.Logger().Debugf("data key=%v is received", data.(*mc.MCRequest).Key)
	xmem.Logger().Debugf("data channel len is %d\n", len(xmem.dataChan))

	request := data.(*mc.MCRequest)

	xmem.dataChan <- request

	xmem.counter_received++

	//accumulate the batchCount and batchSize
	if xmem.batch.accumuBatch(data.(*mc.MCRequest).Size()) {
		xmem.batchReady()
	} else {
		select {
		//not enough for a batch, but the xmem.sendNowCh is signaled
		case <-xmem.sendNowCh:
			xmem.Logger().Debug("Need to send now")
			xmem.batchReady()
		default:
		}
	}
	//raise DataReceived event
	xmem.RaiseEvent(common.DataReceived, data.(*mc.MCRequest), xmem, nil, nil)
	xmem.Logger().Debugf("Xmem %v received %v items\n", xmem.Id(), xmem.counter_received)

	return nil
}

func (xmem *XmemNozzle) processData_async(finch chan bool, waitGrp *sync.WaitGroup) (err error) {
	defer waitGrp.Done()
	for {
		if xmem.IsOpen() {
			select {
			case <-finch:
				goto done
			case batch, ok := <-xmem.batches_ready:
				if batch == nil {
					return nil
				}
				xmem.Logger().Debugf("%v recieved_per_sec=%v, sent_per_sec=%v\n", xmem.Id(),
					float64(xmem.counter_received)/time.Since(xmem.start_time).Seconds(), float64(xmem.counter_received)/time.Since(xmem.start_time).Seconds())
				if !ok {
					return nil
				}
				err = xmem.send_internal(batch)
			}
		}
	}
done:
	return
}

func (xmem *XmemNozzle) processData_batch(finch chan bool, waitGrp *sync.WaitGroup) (err error) {
	defer waitGrp.Done()
	for {
		if xmem.IsOpen() {
			select {
			case <-finch:
				goto done
			case <-xmem.send_allow_ch:
				select {
				case batch, ok := <-xmem.batches_ready:
					xmem.Logger().Debugf("%v Send..., %v batches ready, %v items in queue, count_recieved=%v, count_sent=%v\n", xmem.Id(), len(xmem.batches_ready), len(xmem.dataChan), xmem.counter_received, xmem.counter_sent)
					if !ok {
						return nil
					}
					err = xmem.send_internal(batch)
				default:
					//didn't use the send allowed token, put it back
					select {
					case xmem.send_allow_ch <- true:
					default:
					}
				}
			}
		}
	}
done:
	return
}

func (xmem *XmemNozzle) onExit() {
	//notify the data processing routine
	xmem.sender_finch <- true
	xmem.receiver_finch <- true
	xmem.checker_finch <- true
	xmem.childrenWaitGrp.Wait()

	//cleanup
	pool := base.ConnPoolMgr().GetPool(xmem.getPoolName(xmem.config.connectStr))
	if pool != nil {
		pool.Release(xmem.memClient)
	}

	xmem.memClient = nil
}

func (xmem *XmemNozzle) batchSendWithRetry(batch *xmemBatch, conn io.ReadWriteCloser, numOfRetry int) error {
	var err error
	count := batch.count()

	for i := 0; i < count; i++ {

		item := <-xmem.dataChan
		//blocking
		err, index, reserv_num := xmem.buf.reserveSlot()
		if err != nil {
			return err
		}

		xmem.adjustRequest(item, index)
		item_byte := item.Bytes()

		start_t := time.Now()
		for j := 0; j < numOfRetry; j++ {
			//	fmt.Println(xmem.Id(), "batch transmitting ", string(item.Key))
//			conn.(*net.TCPConn).SetWriteDeadline(time.Now().Add(800 * time.Millisecond))
			_, err = conn.Write(item_byte)
			//	fmt.Println(xmem.Id(), "batch transmitted ", string(item.Key))
			if err == nil {
				break
			}
		}
		xmem.Logger().Infof("send_time=%vs\n", time.Since(start_t).Seconds())

		if err == nil {
			err = xmem.buf.enSlot(index, item, reserv_num)
		}

		if err != nil {
			fmt.Println("Failed to send. err=%v\n", err)
			xmem.buf.cancelReservation(index, reserv_num)
		}

	}

	if conn == nil {
		panic("lost connection")
	}

	//log the data
	return err
}

//func (xmem *XmemNozzle) send() error {
//	var err error
//
//	//get the batch to process
//	if xmem.config.mode == Batch_XMEM {
//		<-xmem.send_allow_ch
//		select {
//		case batch, ok := <-xmem.batches_ready:
//			xmem.Logger().Debugf("%v Send..., %v batches ready, %v items in queue, count_recieved=%v, count_sent=%v\n", xmem.Id(), len(xmem.batches_ready), len(xmem.dataChan), xmem.counter_received, xmem.counter_sent)
//			if !ok {
//				return nil
//			}
//			err = xmem.send_internal(batch)
//		default:
//			//didn't use the send allowed token, put it back
//			select {
//			case xmem.send_allow_ch <- true:
//			default:
//			}
//		}
//	} else {
//		//async mode
//		select {
//		case batch, ok := <-xmem.batches_ready:
//			if batch == nil {
//				return nil
//			}
//			xmem.Logger().Debugf("%v recieved_per_sec=%v, sent_per_sec=%v\n", xmem.Id(),
//				float64(xmem.counter_received)/time.Since(xmem.start_time).Seconds(), float64(xmem.counter_received)/time.Since(xmem.start_time).Seconds())
//			if !ok {
//				return nil
//			}
//			err = xmem.send_internal(batch)
//
//		}
//
//	}
//
//	return err
//}

func (xmem *XmemNozzle) send_internal(batch *xmemBatch) error {
	var err error
	count := batch.count()

	xmem.Logger().Infof("Send batch count=%d\n", count)

	xmem.counter_sent = xmem.counter_sent + count
	xmem.Logger().Debugf("So far, xmem %v processed %d items", xmem.Id(), xmem.counter_sent)

	//get the raw connection
	conn := xmem.memClient.Hijack()
	defer func() {
		//		xmem.memClient, err = mcc.Wrap(conn)

		if err != nil || xmem.memClient == nil {
			xmem.Logger().Errorf("Connection lost, recover....")
			conn.Close()
			err = xmem.recoverFromConnLost()
			if err != nil {
				otherInfo := utils.WrapError(err)
				xmem.RaiseEvent(common.ErrorEncountered, nil, xmem, nil, otherInfo)
			}
		}
	}()

	//batch send
	err = xmem.batchSendWithRetry(batch, conn, xmem.config.maxRetry)
	return err
}

func (xmem *XmemNozzle) sendSingleWithRetry(adjustRequest bool, item *mc.MCRequest, numOfRetry int, index uint16) (err error) {
	for i := 0; i < numOfRetry; i++ {
		err = xmem.sendSingle(adjustRequest, item, index)
		if err == nil {
			break
		}
	}

	return err
}

func (xmem *XmemNozzle) sendSingle(adjustRequest bool, item *mc.MCRequest, index uint16) error {
	if adjustRequest {
		xmem.adjustRequest(item, index)
		xmem.Logger().Debugf("key=%v\n", item.Key)
		xmem.Logger().Debugf("opcode=%v\n", item.Opcode)
	}
	if xmem.memClient == nil {
		xmem.Logger().Errorf("%v, Connection lost, recover....", xmem.Id())
		err := xmem.recoverFromConnLost()
		if err != nil {
			otherInfo := utils.WrapError(err)
			xmem.RaiseEvent(common.ErrorEncountered, nil, xmem, nil, otherInfo)
		}
	}
	bytes := item.Bytes()
	conn := xmem.memClient.Hijack()

//	fmt.Printf("%v %v transmitting %s %x pos =%v\n", time.Now(), xmem.Id(), string(item.Key), item.Opaque, index)
//	conn.(*net.TCPConn).SetWriteDeadline(time.Now().Add(800 * time.Millisecond))
	_, err := conn.Write(bytes)
//	fmt.Printf("%v %v transmited %s %x pos=%v\n", time.Now(), xmem.Id(), string(item.Key), item.Opaque, index)

	if err != nil {
		if neterr, ok := err.(net.Error); ok && (neterr.Timeout() || neterr.Temporary()) {
			xmem.Logger().Errorf("%v sendSingle: transmit error: %s\n", xmem.Id(), fmt.Sprint(err))
			return err
		} else {
			panic(err)
		}
	}

	return nil
}

//TODO: who will release the pool? maybe it should be replication manager
//
func (xmem *XmemNozzle) initializeConnection() (err error) {
	xmem.Logger().Debugf("xmem.config= %v", xmem.config.connectStr)
	xmem.Logger().Debugf("poolName=%v", xmem.getPoolName(xmem.config.connectStr))
	pool, err := base.ConnPoolMgr().GetOrCreatePool(xmem.getPoolName(xmem.config.connectStr), xmem.config.connectStr, xmem.config.bucketName, xmem.config.password, base.DefaultConnectionSize)
	if err == nil {
		xmem.memClient, err = pool.Get()
	}
	return err
}

func (xmem *XmemNozzle) getPoolName(connectionStr string) string {
	return "Couch_Xmem_" + connectionStr
}

func (xmem *XmemNozzle) initNewBatch() {
	xmem.Logger().Info("init a new batch")
	xmem.batch = newXmemBatch(xmem.config.maxCount, xmem.config.maxSize, xmem.Logger())
}

func (xmem *XmemNozzle) initialize(settings map[string]interface{}) error {
	err := xmem.config.initializeConfig(settings)
	xmem.dataChan = make(chan *mc.MCRequest, xmem.config.maxCount*100)
	xmem.sendNowCh = make(chan bool, 1)

	//enable send
	xmem.send_allow_ch <- true

	//init a new batch
	xmem.initNewBatch()

	xmem.batch_move_ch = make(chan bool, 1)
	xmem.batch_move_ch <- true

	//initialize buffer
	if xmem.config.mode == Batch_XMEM {
		xmem.buf = newReqBuffer(uint16(xmem.config.maxCount*200), xmem.send_allow_ch, xmem.Logger())
	} else {
		//no send batch control
		close(xmem.send_allow_ch)
		xmem.buf = newReqBuffer(uint16(xmem.config.maxCount*200), nil, xmem.Logger())
	}

	xmem.receiver_finch = make(chan bool, 1)
	xmem.checker_finch = make(chan bool, 1)

	xmem.Logger().Debug("About to start initializing connection")
	if err == nil {
		err = xmem.initializeConnection()
	}

	xmem.Logger().Debug("Initialization is done")

	return err
}

func (xmem *XmemNozzle) receiveResponse(client *mcc.Client, finch chan bool, waitGrp *sync.WaitGroup) {
	defer waitGrp.Done()

	var count = 0
	for {
		response, err := client.Receive()
		count++
		if err != io.EOF {
			pos := xmem.getPosFromOpaque(response.Opaque)
			if err != nil && isRealError(response.Status) {
				if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
					xmem.Logger().Errorf("read time out, exiting..., err=%v\n", neterr)
					goto done
				}
				xmem.Logger().Errorf("%v pos=%d, Received error = %v in response, err = %v, response=%v\n", xmem.Id(), pos, response.Status.String(), err, response.Bytes())
				_, err = xmem.buf.modSlot(pos, xmem.resend)
			} else {
				//raiseEvent
				req, _ := xmem.buf.slot(pos)
				if req != nil && req.Opaque == response.Opaque {
					xmem.Logger().Debugf("%v received opaque=%v\n", xmem.Id(), req.Opaque)
					xmem.RaiseEvent(common.DataSent, req, xmem, nil, nil)
					//empty the slot in the buffer
					if xmem.buf.evictSlot(pos) != nil {
						xmem.Logger().Errorf("Failed to evict slot %d\n", pos)
					}
				} else {
					//probably redundant response we got from our resend
					//log and ignore
					if req == nil {
						//xmem.Logger().Errorf("%v Received response for req is nil, pos=%v, response.Opaque=%x again, %v empty slots \n", xmem.Id(), pos, response.Opaque, len(xmem.buf.empty_slots_pos))
					} else {
						//xmem.Logger().Errorf("%v Received response for req.Opaque=%x, pos=%v, response.Opaque=%x again\n", xmem.Id(), req.Opaque, pos, response.Opaque)
					}
				}
			}
		} else {
			xmem.Logger().Errorf("Quit receiveResponse. err=%v\n", err)
			goto done
		}
	}

done:
}

func isRealError(resp_status mc.Status) bool {
	switch resp_status {
	case mc.KEY_ENOENT, mc.KEY_EEXISTS, mc.NOT_STORED, mc.Status(0x87):
		return false
	default:
		return true
	}
}
func (xmem *XmemNozzle) check(finch chan bool, waitGrp *sync.WaitGroup) {
	var count uint64
	ticker := time.Tick(xmem.config.batchtimeout)
	for {
		select {
		case <-finch:
			goto done
		case <-ticker:
			count++
			if math.Mod(float64(count), float64(15)) < 2 {
				xmem.Logger().Debugf("%v checking timeout. %v unsent\n", xmem.Id(), len(xmem.dataChan))
			}
			size := xmem.buf.bufferSize()
			timeoutCheckFunc := xmem.checkTimeout
			for i := 0; i < int(size); i++ {
				xmem.buf.modSlot(uint16(i), timeoutCheckFunc)
			}

			if xmem.isCurrentBatchExpiring() {
				xmem.batchReady()
			}

		}
	}
done:
	xmem.Logger().Debug("Xmem checking routine exits")
	waitGrp.Done()
}

func (xmem *XmemNozzle) checkTimeout(req *bufferedMCRequest, pos uint16) bool {
	duration := xmem.timeoutDuration(req.num_of_retry)
	if time.Since(req.sent_time) > duration {
		xmem.Logger().Debugf("key=%v, numOfRetry=%v, timeout duration=%v\n", string(req.req.Key), req.num_of_retry, duration)
		modified := xmem.resend(req, pos)
		return modified
	}
	return false
}

func (xmem *XmemNozzle) timeoutDuration(numofRetry int) time.Duration {
	duration := 0 * time.Millisecond
	for i := 0; i <= numofRetry; i++ {
		duration = duration + time.Duration(i+1)*xmem.config.batchtimeout
	}
	return duration
}

func (xmem *XmemNozzle) resend(req *bufferedMCRequest, pos uint16) bool {
	if req.num_of_retry < xmem.config.maxRetry-1 {
		xmem.Logger().Debug("Retry sending ....")
		err := xmem.sendSingle(false, req.req, pos)

		if err != nil {
			req.err = err
		} else {
			req.num_of_retry = req.num_of_retry + 1
		}
		return true
	} else {
		//raise error
		xmem.Logger().Errorf("%v Max number of retry has been reached for key=%s, wait_time=%v\n",
			xmem.Id(), req.req.Key, time.Since(req.sent_time))
		err := errors.New("Max number of retry has reached")
		otherInfo := utils.WrapError(err)
		xmem.RaiseEvent(common.ErrorEncountered, req.req, xmem, nil, otherInfo)
		xmem.buf.evictSlot(pos)
	}
	return false
}

func (xmem *XmemNozzle) adjustRequest(mc_req *mc.MCRequest, index uint16) {
	mc_req.Opcode = xmem.encodeOpCode(mc_req.Opcode)
	mc_req.Cas = 0
	mc_req.Opaque = xmem.getOpaque(index, xmem.buf.sequences[int(index)])
	mc_req.Extras = []byte{0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0}
	binary.BigEndian.PutUint64(mc_req.Extras, uint64(0)<<32|uint64(0))

}

func (xmem *XmemNozzle) getOpaque(index, sequence uint16) uint32 {
	result := uint32(sequence)<<16 + uint32(index)
	xmem.Logger().Debugf("uint32(sequence)<<16 = %v", uint32(sequence)<<16)
	xmem.Logger().Debugf("index=%x, sequence=%x, opaque=%x\n", index, sequence, result)
	return result
}

func (xmem *XmemNozzle) getPosFromOpaque(opaque uint32) uint16 {
	result := uint16(0x0000FFFF & opaque)
	xmem.Logger().Debugf("opaque=%x, index=%v\n", opaque, result)
	return result
}

func (xmem *XmemNozzle) encodeOpCode(code mc.CommandCode) mc.CommandCode {
	if code == mc.UPR_MUTATION || code == mc.TAP_MUTATION {
		return SET_WITH_META
	} else if code == mc.TAP_DELETE || code == mc.UPR_DELETION {
		return DELETE_WITH_META
	}
	return code
}

func (xmem *XmemNozzle) isCurrentBatchExpiring() bool {
	if xmem.batch.count() >= 1 && time.Since(xmem.batch.start_time) > xmem.config.batchExpirationTime {
		xmem.Logger().Debugf("%v the current batch count=%v is expiring, ready batch=%v, time passed %v \n", xmem.Id(), xmem.batch.count(), len(xmem.batches_ready), time.Since(xmem.batch.start_time))
		return true
	}
	return false

}

func (xmem *XmemNozzle) StatusSummary() string {
	return fmt.Sprintf("Xmem %v received %v items, sent %v items", xmem.Id(), xmem.counter_received, xmem.counter_sent)
}

func (xmem *XmemNozzle) handleGeneralError(err error) {
	otherInfo := utils.WrapError(err)
	xmem.RaiseEvent(common.ErrorEncountered, nil, xmem, nil, otherInfo)
}

func (xmem *XmemNozzle) recoverFromConnLost() (err error) {
	//failed to recycle the connection
	//reinitialize the connection
	err = xmem.initializeConnection()
	if err != nil {
		xmem.Logger().Errorf("failed to recycle the connection, err=%v", err)
		return
	}

	xmem.childrenWaitGrp.Add(1)
	go xmem.receiveResponse(xmem.memClient, xmem.receiver_finch, &xmem.childrenWaitGrp)
	return
}
