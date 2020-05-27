// Copyright (c) 2013-2020 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.

package metadata

import (
	"encoding/json"
	"fmt"
	"github.com/couchbase/goxdcr/base"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBackfillReplMarshal(t *testing.T) {
	fmt.Println("============== Test case start: TestBackfillReplMarshal =================")
	assert := assert.New(t)

	namespaceMapping := make(CollectionNamespaceMapping)
	defaultNamespace := &base.CollectionNamespace{base.DefaultScopeCollectionName, base.DefaultScopeCollectionName}
	namespaceMapping.AddSingleMapping(defaultNamespace, defaultNamespace)

	manifestsIdPair := base.CollectionsManifestIdPair{0, 0}
	ts0 := &BackfillVBTimestamps{
		StartingTimestamp: &base.VBTimestamp{0, 0, 5, 10, 10, manifestsIdPair},
		EndingTimestamp:   &base.VBTimestamp{0, 0, 5000, 500, 500, manifestsIdPair},
	}

	vb0Task0 := NewBackfillTask(ts0, []CollectionNamespaceMapping{namespaceMapping})

	ts1 := &BackfillVBTimestamps{
		StartingTimestamp: &base.VBTimestamp{0, 0, 5005, 10, 10, manifestsIdPair},
		EndingTimestamp:   &base.VBTimestamp{0, 0, 15005, 500, 500, manifestsIdPair},
	}
	vb0Task1 := NewBackfillTask(ts1, []CollectionNamespaceMapping{namespaceMapping})
	_, err := json.Marshal(vb0Task0)
	assert.Nil(err)

	var vb0Tasks BackfillTasks
	vb0Tasks = append(vb0Tasks, vb0Task0)
	vb0Tasks = append(vb0Tasks, vb0Task1)

	ts2 := &BackfillVBTimestamps{
		StartingTimestamp: &base.VBTimestamp{1, 0, 5, 10, 10, manifestsIdPair},
		EndingTimestamp:   &base.VBTimestamp{1, 0, 5000, 500, 500, manifestsIdPair},
	}
	vb1Task0 := NewBackfillTask(ts2, []CollectionNamespaceMapping{namespaceMapping})

	var vb1Tasks BackfillTasks
	vb1Tasks = append(vb1Tasks, vb1Task0)

	_, err = json.Marshal(vb1Tasks)
	assert.Nil(err)

	vbTasksMap := make(map[uint16]*BackfillTasks)
	vbTasksMap[0] = &vb0Tasks
	vbTasksMap[1] = &vb1Tasks

	testId := "testId"
	testInternalId := "testInternalId"
	testSpec := &BackfillReplicationSpec{
		Id:         testId,
		InternalId: testInternalId,
		VBTasksMap: vbTasksMap,
	}

	marshalledSpec, err := json.Marshal(testSpec)
	assert.Nil(err)

	checkSpec := &BackfillReplicationSpec{
		VBTasksMap: make(map[uint16]*BackfillTasks),
	}
	err = json.Unmarshal(marshalledSpec, &checkSpec)
	assert.Nil(err)
	assert.True(checkSpec.SameAs(testSpec))

	assert.Equal(2, len(*(testSpec.VBTasksMap[0])))
	assert.Equal(1, len(*(testSpec.VBTasksMap[1])))
	assert.Nil(testSpec.VBTasksMap[2])

	// Test append
	vb0TasksClone := vb0Tasks.Clone()
	ts0_append := &BackfillVBTimestamps{
		StartingTimestamp: &base.VBTimestamp{0, 0, 5, 10, 10, manifestsIdPair},
		EndingTimestamp:   &base.VBTimestamp{0, 0, 7000, 700, 700, manifestsIdPair},
	}
	vb0AppendTask0 := NewBackfillTask(ts0_append, []CollectionNamespaceMapping{namespaceMapping})
	ts2_append := &BackfillVBTimestamps{
		StartingTimestamp: &base.VBTimestamp{1, 0, 5, 10, 10, manifestsIdPair},
		EndingTimestamp:   &base.VBTimestamp{1, 0, 10000, 1000, 1000, manifestsIdPair},
	}
	vb1AppendTask0 := NewBackfillTask(ts2_append, []CollectionNamespaceMapping{namespaceMapping})
	var vb0TasksAppend BackfillTasks
	vb0TasksAppend = append(vb0TasksAppend, vb0AppendTask0)
	var vb1TasksAppend BackfillTasks
	vb1TasksAppend = append(vb1TasksAppend, vb1AppendTask0)
	appendTasksMap := make(map[uint16]*BackfillTasks)
	appendTasksMap[0] = &vb0TasksAppend
	appendTasksMap[1] = &vb1TasksAppend
	// Add an create for good measure
	appendTasksMap[2] = &vb0TasksClone

	testSpec.AppendTasks(appendTasksMap)

	assert.Equal(3, len(*(testSpec.VBTasksMap[0])))
	assert.Equal(2, len(*(testSpec.VBTasksMap[1])))
	assert.Equal(2, len(*(testSpec.VBTasksMap[2])))

	testSpec.VBTasksMap.MarkOneVBTaskDone(0)
	assert.Equal(2, len(*(testSpec.VBTasksMap[0])))
	fmt.Println("============== Test case end: TestBackfillReplMarshal =================")
}

func TestVBTimestampAccomodate(t *testing.T) {
	assert := assert.New(t)
	fmt.Println("============== Test case start: TestVBTimestampAccomodate =================")

	origBeg := &base.VBTimestamp{Seqno: 10}
	origEnd := &base.VBTimestamp{Seqno: 50}
	origTimestamp := &BackfillVBTimestamps{origBeg, origEnd}

	fullyAccomodateTs := &BackfillVBTimestamps{origBeg, origEnd}
	fullyAccomodated, unableToMerge, smallerOob, largerOob := origTimestamp.Accomodate(fullyAccomodateTs)

	assert.True(fullyAccomodated)
	assert.Nil(smallerOob.StartingTimestamp)
	assert.Nil(smallerOob.EndingTimestamp)
	assert.Nil(largerOob.StartingTimestamp)
	assert.Nil(largerOob.EndingTimestamp)
	assert.False(unableToMerge)

	extraBeg := &base.VBTimestamp{Seqno: 0}
	trailingTail := &BackfillVBTimestamps{extraBeg, origEnd}
	fullyAccomodated, unableToMerge, smallerOob, largerOob = origTimestamp.Accomodate(trailingTail)
	assert.False(fullyAccomodated)
	assert.Nil(largerOob.StartingTimestamp)
	assert.Nil(largerOob.EndingTimestamp)
	assert.Equal(uint64(0), smallerOob.StartingTimestamp.Seqno)
	assert.Equal(origBeg.Seqno, smallerOob.EndingTimestamp.Seqno)
	assert.False(unableToMerge)

	extraEnd := &base.VBTimestamp{Seqno: 60}
	fullyTrail := &BackfillVBTimestamps{extraBeg, extraEnd}
	fullyAccomodated, unableToMerge, smallerOob, largerOob = origTimestamp.Accomodate(fullyTrail)
	assert.False(fullyAccomodated)
	assert.Equal(uint64(0), smallerOob.StartingTimestamp.Seqno)
	assert.Equal(origBeg.Seqno, smallerOob.EndingTimestamp.Seqno)
	assert.Equal(origEnd.Seqno, largerOob.StartingTimestamp.Seqno)
	assert.Equal(extraEnd.Seqno, largerOob.EndingTimestamp.Seqno)
	assert.False(unableToMerge)

	fmt.Println("============== Test case end: TestVBTimestampAccomodate =================")
}

func TestMergeTask(t *testing.T) {
	assert := assert.New(t)
	fmt.Println("============== Test case start: TestMergeTask =================")
	namespaceMapping := make(CollectionNamespaceMapping)
	defaultNamespace := &base.CollectionNamespace{base.DefaultScopeCollectionName, base.DefaultScopeCollectionName}
	namespaceMapping.AddSingleMapping(defaultNamespace, defaultNamespace)

	namespaceMapping2 := make(CollectionNamespaceMapping)
	namespace2 := &base.CollectionNamespace{"dummy", "dummy"}
	namespaceMapping2.AddSingleMapping(namespace2, namespace2)

	manifestsIdPair := base.CollectionsManifestIdPair{0, 0}
	ts0 := &BackfillVBTimestamps{
		StartingTimestamp: &base.VBTimestamp{0, 0, 5, 10, 10, manifestsIdPair},
		EndingTimestamp:   &base.VBTimestamp{0, 0, 5000, 500, 500, manifestsIdPair},
	}
	task0 := NewBackfillTask(ts0, []CollectionNamespaceMapping{namespaceMapping})

	ts1 := &BackfillVBTimestamps{
		StartingTimestamp: &base.VBTimestamp{0, 0, 5, 10, 10, manifestsIdPair},
		EndingTimestamp:   &base.VBTimestamp{0, 0, 100, 500, 500, manifestsIdPair},
	}
	task1 := NewBackfillTask(ts1, []CollectionNamespaceMapping{namespaceMapping})

	canFullyMerge, unableToMerge, subTask1, subTask2 := task0.MergeIncomingTask(task1)
	assert.True(canFullyMerge)
	assert.False(unableToMerge)
	assert.Nil(subTask1)
	assert.Nil(subTask2)

	ts2 := &BackfillVBTimestamps{
		StartingTimestamp: &base.VBTimestamp{0, 0, 0, 10, 10, manifestsIdPair},
		EndingTimestamp:   &base.VBTimestamp{0, 0, 100, 500, 500, manifestsIdPair},
	}
	task2 := NewBackfillTask(ts2, []CollectionNamespaceMapping{namespaceMapping2})
	canFullyMerge, unableToMerge, subTask1, subTask2 = task0.MergeIncomingTask(task2)
	assert.False(canFullyMerge)
	assert.False(unableToMerge)
	assert.NotNil(subTask1)
	assert.Equal(uint64(0), subTask1.Timestamps.StartingTimestamp.Seqno)
	assert.Equal(uint64(5), subTask1.Timestamps.EndingTimestamp.Seqno)
	assert.Nil(subTask2)
	assert.Equal(2, len(task0.RequestedCollections()))

	ts3 := &BackfillVBTimestamps{
		StartingTimestamp: &base.VBTimestamp{0, 0, 5001, 10, 10, manifestsIdPair},
		EndingTimestamp:   &base.VBTimestamp{0, 0, 10000, 500, 500, manifestsIdPair},
	}
	task3 := NewBackfillTask(ts3, []CollectionNamespaceMapping{namespaceMapping})
	canFullyMerge, unableToMerge, subTask1, subTask2 = task0.MergeIncomingTask(task3)
	assert.False(canFullyMerge)
	assert.True(unableToMerge)
	assert.Nil(subTask1)
	assert.Nil(subTask2)
	fmt.Println("============== Test case end: TestMergeTask =================")
}

func TestMergeTasks(t *testing.T) {
	assert := assert.New(t)
	fmt.Println("============== Test case start: TestMergeTasks =================")

	namespaceMapping := make(CollectionNamespaceMapping)
	defaultNamespace := &base.CollectionNamespace{base.DefaultScopeCollectionName, base.DefaultScopeCollectionName}
	namespaceMapping.AddSingleMapping(defaultNamespace, defaultNamespace)

	manifestsIdPair := base.CollectionsManifestIdPair{0, 0}
	ts0 := &BackfillVBTimestamps{
		StartingTimestamp: &base.VBTimestamp{0, 0, 5, 10, 10, manifestsIdPair},
		EndingTimestamp:   &base.VBTimestamp{0, 0, 5000, 500, 500, manifestsIdPair},
	}

	vb0Task0 := NewBackfillTask(ts0, []CollectionNamespaceMapping{namespaceMapping})

	ts1 := &BackfillVBTimestamps{
		StartingTimestamp: &base.VBTimestamp{0, 0, 5005, 10, 10, manifestsIdPair},
		EndingTimestamp:   &base.VBTimestamp{0, 0, 15005, 500, 500, manifestsIdPair},
	}
	vb0Task1 := NewBackfillTask(ts1, []CollectionNamespaceMapping{namespaceMapping})

	var totalTasks BackfillTasks
	totalTasks = append(totalTasks, vb0Task0)
	totalTasks = append(totalTasks, vb0Task1)

	// Now try to merge a task that should overlaps twice
	ts2 := &BackfillVBTimestamps{
		StartingTimestamp: &base.VBTimestamp{0, 0, 0, 10, 10, manifestsIdPair},
		EndingTimestamp:   &base.VBTimestamp{0, 0, 20000, 500, 500, manifestsIdPair},
	}
	vb0Task2 := NewBackfillTask(ts2, []CollectionNamespaceMapping{namespaceMapping})

	var unmergableTasks BackfillTasks
	assert.Equal(0, len(unmergableTasks))
	totalTasks.MergeIncomingTaskIntoTasks(vb0Task2, &unmergableTasks)
	assert.Equal(3, len(unmergableTasks))

	assert.True(unmergableTasks.containsStartEndRange(0, 5))
	assert.True(unmergableTasks.containsStartEndRange(5000, 5005))
	assert.True(unmergableTasks.containsStartEndRange(15005, 20000))
	fmt.Println("============== Test case end: TestMergeTasks =================")
}

func TestMergeTasksIntoSpec(t *testing.T) {
	fmt.Println("============== Test case start: TestMergeTasksIntoSpec =================")
	assert := assert.New(t)

	namespaceMapping := make(CollectionNamespaceMapping)
	defaultNamespace := &base.CollectionNamespace{base.DefaultScopeCollectionName, base.DefaultScopeCollectionName}
	namespaceMapping.AddSingleMapping(defaultNamespace, defaultNamespace)

	manifestsIdPair := base.CollectionsManifestIdPair{0, 0}
	ts0 := &BackfillVBTimestamps{
		StartingTimestamp: &base.VBTimestamp{0, 0, 5, 10, 10, manifestsIdPair},
		EndingTimestamp:   &base.VBTimestamp{0, 0, 5000, 500, 500, manifestsIdPair},
	}

	vb0Task0 := NewBackfillTask(ts0, []CollectionNamespaceMapping{namespaceMapping})

	ts1 := &BackfillVBTimestamps{
		StartingTimestamp: &base.VBTimestamp{0, 0, 5005, 10, 10, manifestsIdPair},
		EndingTimestamp:   &base.VBTimestamp{0, 0, 15005, 500, 500, manifestsIdPair},
	}
	vb0Task1 := NewBackfillTask(ts1, []CollectionNamespaceMapping{namespaceMapping})

	var vb0Tasks BackfillTasks
	vb0Tasks = append(vb0Tasks, vb0Task0)
	vb0Tasks = append(vb0Tasks, vb0Task1)

	vbTasksMap := make(map[uint16]*BackfillTasks)
	vbTasksMap[0] = &vb0Tasks

	testId := "testId"
	testInternalId := "testInternalId"
	testSpec := &BackfillReplicationSpec{
		Id:         testId,
		InternalId: testInternalId,
		VBTasksMap: vbTasksMap,
	}

	// This test will also ensure that the first one is skipped to simulate that
	// the first task is undergoing backfill
	// So when the new task comes in, it needs to re-attempt to refill from seqno 5
	newTs := &BackfillVBTimestamps{
		StartingTimestamp: &base.VBTimestamp{0, 0, 5, 10, 10, manifestsIdPair},
		EndingTimestamp:   &base.VBTimestamp{0, 0, 20000, 500, 500, manifestsIdPair},
	}
	newVb0Task := NewBackfillTask(newTs, []CollectionNamespaceMapping{namespaceMapping})
	var newVb0Tasks BackfillTasks
	newVb0Tasks = append(newVb0Tasks, newVb0Task)
	newVbTaskMap := make(map[uint16]*BackfillTasks)
	newVbTaskMap[0] = &newVb0Tasks

	testSpec.MergeNewTasks(newVbTaskMap, true /*skipFirst*/)

	assert.Equal(4, len(*(testSpec.VBTasksMap[0])))
	// This is the first ongoing backfill
	assert.True(testSpec.VBTasksMap[0].containsStartEndRange(5, 5000))
	// Taking the first one away, we are left with
	assert.True(testSpec.VBTasksMap[0].containsStartEndRange(5, 5005))
	assert.True(testSpec.VBTasksMap[0].containsStartEndRange(5005, 15005))
	assert.True(testSpec.VBTasksMap[0].containsStartEndRange(15005, 20000))
	fmt.Println("============== Test case end: TestMergeTasksIntoSpec =================")
}
