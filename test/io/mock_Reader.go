// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MockReader is an autogenerated mock type for the Reader type
type MockReader struct {
	mock.Mock
}

type MockReader_Expecter struct {
	mock *mock.Mock
}

func (_m *MockReader) EXPECT() *MockReader_Expecter {
	return &MockReader_Expecter{mock: &_m.Mock}
}

// Read provides a mock function with given fields: p
func (_m *MockReader) Read(p []byte) (int, error) {
	ret := _m.Called(p)

	if len(ret) == 0 {
		panic("no return value specified for Read")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (int, error)); ok {
		return rf(p)
	}
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockReader_Read_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Read'
type MockReader_Read_Call struct {
	*mock.Call
}

// Read is a helper method to define mock.On call
//   - p []byte
func (_e *MockReader_Expecter) Read(p interface{}) *MockReader_Read_Call {
	return &MockReader_Read_Call{Call: _e.mock.On("Read", p)}
}

func (_c *MockReader_Read_Call) Run(run func(p []byte)) *MockReader_Read_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *MockReader_Read_Call) Return(n int, err error) *MockReader_Read_Call {
	_c.Call.Return(n, err)
	return _c
}

func (_c *MockReader_Read_Call) RunAndReturn(run func([]byte) (int, error)) *MockReader_Read_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockReader creates a new instance of MockReader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockReader {
	mock := &MockReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
