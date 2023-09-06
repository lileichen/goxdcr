// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UIWarnings is an autogenerated mock type for the UIWarnings type
type UIWarnings struct {
	mock.Mock
}

type UIWarnings_Expecter struct {
	mock *mock.Mock
}

func (_m *UIWarnings) EXPECT() *UIWarnings_Expecter {
	return &UIWarnings_Expecter{mock: &_m.Mock}
}

// AddWarning provides a mock function with given fields: key, val
func (_m *UIWarnings) AddWarning(key string, val string) {
	_m.Called(key, val)
}

// UIWarnings_AddWarning_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddWarning'
type UIWarnings_AddWarning_Call struct {
	*mock.Call
}

// AddWarning is a helper method to define mock.On call
//   - key string
//   - val string
func (_e *UIWarnings_Expecter) AddWarning(key interface{}, val interface{}) *UIWarnings_AddWarning_Call {
	return &UIWarnings_AddWarning_Call{Call: _e.mock.On("AddWarning", key, val)}
}

func (_c *UIWarnings_AddWarning_Call) Run(run func(key string, val string)) *UIWarnings_AddWarning_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *UIWarnings_AddWarning_Call) Return() *UIWarnings_AddWarning_Call {
	_c.Call.Return()
	return _c
}

func (_c *UIWarnings_AddWarning_Call) RunAndReturn(run func(string, string)) *UIWarnings_AddWarning_Call {
	_c.Call.Return(run)
	return _c
}

// AppendGeneric provides a mock function with given fields: warning
func (_m *UIWarnings) AppendGeneric(warning string) {
	_m.Called(warning)
}

// UIWarnings_AppendGeneric_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AppendGeneric'
type UIWarnings_AppendGeneric_Call struct {
	*mock.Call
}

// AppendGeneric is a helper method to define mock.On call
//   - warning string
func (_e *UIWarnings_Expecter) AppendGeneric(warning interface{}) *UIWarnings_AppendGeneric_Call {
	return &UIWarnings_AppendGeneric_Call{Call: _e.mock.On("AppendGeneric", warning)}
}

func (_c *UIWarnings_AppendGeneric_Call) Run(run func(warning string)) *UIWarnings_AppendGeneric_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *UIWarnings_AppendGeneric_Call) Return() *UIWarnings_AppendGeneric_Call {
	_c.Call.Return()
	return _c
}

func (_c *UIWarnings_AppendGeneric_Call) RunAndReturn(run func(string)) *UIWarnings_AppendGeneric_Call {
	_c.Call.Return(run)
	return _c
}

// GetFieldWarningsOnly provides a mock function with given fields:
func (_m *UIWarnings) GetFieldWarningsOnly() map[string]interface{} {
	ret := _m.Called()

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func() map[string]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

// UIWarnings_GetFieldWarningsOnly_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFieldWarningsOnly'
type UIWarnings_GetFieldWarningsOnly_Call struct {
	*mock.Call
}

// GetFieldWarningsOnly is a helper method to define mock.On call
func (_e *UIWarnings_Expecter) GetFieldWarningsOnly() *UIWarnings_GetFieldWarningsOnly_Call {
	return &UIWarnings_GetFieldWarningsOnly_Call{Call: _e.mock.On("GetFieldWarningsOnly")}
}

func (_c *UIWarnings_GetFieldWarningsOnly_Call) Run(run func()) *UIWarnings_GetFieldWarningsOnly_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UIWarnings_GetFieldWarningsOnly_Call) Return(_a0 map[string]interface{}) *UIWarnings_GetFieldWarningsOnly_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UIWarnings_GetFieldWarningsOnly_Call) RunAndReturn(run func() map[string]interface{}) *UIWarnings_GetFieldWarningsOnly_Call {
	_c.Call.Return(run)
	return _c
}

// GetSuccessfulWarningStrings provides a mock function with given fields:
func (_m *UIWarnings) GetSuccessfulWarningStrings() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// UIWarnings_GetSuccessfulWarningStrings_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSuccessfulWarningStrings'
type UIWarnings_GetSuccessfulWarningStrings_Call struct {
	*mock.Call
}

// GetSuccessfulWarningStrings is a helper method to define mock.On call
func (_e *UIWarnings_Expecter) GetSuccessfulWarningStrings() *UIWarnings_GetSuccessfulWarningStrings_Call {
	return &UIWarnings_GetSuccessfulWarningStrings_Call{Call: _e.mock.On("GetSuccessfulWarningStrings")}
}

func (_c *UIWarnings_GetSuccessfulWarningStrings_Call) Run(run func()) *UIWarnings_GetSuccessfulWarningStrings_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UIWarnings_GetSuccessfulWarningStrings_Call) Return(_a0 []string) *UIWarnings_GetSuccessfulWarningStrings_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UIWarnings_GetSuccessfulWarningStrings_Call) RunAndReturn(run func() []string) *UIWarnings_GetSuccessfulWarningStrings_Call {
	_c.Call.Return(run)
	return _c
}

// Len provides a mock function with given fields:
func (_m *UIWarnings) Len() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// UIWarnings_Len_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Len'
type UIWarnings_Len_Call struct {
	*mock.Call
}

// Len is a helper method to define mock.On call
func (_e *UIWarnings_Expecter) Len() *UIWarnings_Len_Call {
	return &UIWarnings_Len_Call{Call: _e.mock.On("Len")}
}

func (_c *UIWarnings_Len_Call) Run(run func()) *UIWarnings_Len_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UIWarnings_Len_Call) Return(_a0 int) *UIWarnings_Len_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UIWarnings_Len_Call) RunAndReturn(run func() int) *UIWarnings_Len_Call {
	_c.Call.Return(run)
	return _c
}

// String provides a mock function with given fields:
func (_m *UIWarnings) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// UIWarnings_String_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'String'
type UIWarnings_String_Call struct {
	*mock.Call
}

// String is a helper method to define mock.On call
func (_e *UIWarnings_Expecter) String() *UIWarnings_String_Call {
	return &UIWarnings_String_Call{Call: _e.mock.On("String")}
}

func (_c *UIWarnings_String_Call) Run(run func()) *UIWarnings_String_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UIWarnings_String_Call) Return(_a0 string) *UIWarnings_String_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UIWarnings_String_Call) RunAndReturn(run func() string) *UIWarnings_String_Call {
	_c.Call.Return(run)
	return _c
}

// NewUIWarnings creates a new instance of UIWarnings. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUIWarnings(t interface {
	mock.TestingT
	Cleanup(func())
}) *UIWarnings {
	mock := &UIWarnings{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}