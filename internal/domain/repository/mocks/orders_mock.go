// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/angryscorp/gophermart/internal/domain/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// OrdersMock is an autogenerated mock type for the Orders type
type OrdersMock struct {
	mock.Mock
}

type OrdersMock_Expecter struct {
	mock *mock.Mock
}

func (_m *OrdersMock) EXPECT() *OrdersMock_Expecter {
	return &OrdersMock_Expecter{mock: &_m.Mock}
}

// AllOrders provides a mock function with given fields: ctx, userID
func (_m *OrdersMock) AllOrders(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for AllOrders")
	}

	var r0 []model.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) ([]model.Order, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []model.Order); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrdersMock_AllOrders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AllOrders'
type OrdersMock_AllOrders_Call struct {
	*mock.Call
}

// AllOrders is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
func (_e *OrdersMock_Expecter) AllOrders(ctx interface{}, userID interface{}) *OrdersMock_AllOrders_Call {
	return &OrdersMock_AllOrders_Call{Call: _e.mock.On("AllOrders", ctx, userID)}
}

func (_c *OrdersMock_AllOrders_Call) Run(run func(ctx context.Context, userID uuid.UUID)) *OrdersMock_AllOrders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *OrdersMock_AllOrders_Call) Return(_a0 []model.Order, _a1 error) *OrdersMock_AllOrders_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrdersMock_AllOrders_Call) RunAndReturn(run func(context.Context, uuid.UUID) ([]model.Order, error)) *OrdersMock_AllOrders_Call {
	_c.Call.Return(run)
	return _c
}

// CreateOrder provides a mock function with given fields: ctx, order
func (_m *OrdersMock) CreateOrder(ctx context.Context, order model.Order) error {
	ret := _m.Called(ctx, order)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.Order) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OrdersMock_CreateOrder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateOrder'
type OrdersMock_CreateOrder_Call struct {
	*mock.Call
}

// CreateOrder is a helper method to define mock.On call
//   - ctx context.Context
//   - order model.Order
func (_e *OrdersMock_Expecter) CreateOrder(ctx interface{}, order interface{}) *OrdersMock_CreateOrder_Call {
	return &OrdersMock_CreateOrder_Call{Call: _e.mock.On("CreateOrder", ctx, order)}
}

func (_c *OrdersMock_CreateOrder_Call) Run(run func(ctx context.Context, order model.Order)) *OrdersMock_CreateOrder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.Order))
	})
	return _c
}

func (_c *OrdersMock_CreateOrder_Call) Return(_a0 error) *OrdersMock_CreateOrder_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OrdersMock_CreateOrder_Call) RunAndReturn(run func(context.Context, model.Order) error) *OrdersMock_CreateOrder_Call {
	_c.Call.Return(run)
	return _c
}

// OrderInfoForUpdate provides a mock function with given fields: ctx, number
func (_m *OrdersMock) OrderInfoForUpdate(ctx context.Context, number string) (*model.Order, error) {
	ret := _m.Called(ctx, number)

	if len(ret) == 0 {
		panic("no return value specified for OrderInfoForUpdate")
	}

	var r0 *model.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Order, error)); ok {
		return rf(ctx, number)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Order); ok {
		r0 = rf(ctx, number)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, number)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrdersMock_OrderInfoForUpdate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OrderInfoForUpdate'
type OrdersMock_OrderInfoForUpdate_Call struct {
	*mock.Call
}

// OrderInfoForUpdate is a helper method to define mock.On call
//   - ctx context.Context
//   - number string
func (_e *OrdersMock_Expecter) OrderInfoForUpdate(ctx interface{}, number interface{}) *OrdersMock_OrderInfoForUpdate_Call {
	return &OrdersMock_OrderInfoForUpdate_Call{Call: _e.mock.On("OrderInfoForUpdate", ctx, number)}
}

func (_c *OrdersMock_OrderInfoForUpdate_Call) Run(run func(ctx context.Context, number string)) *OrdersMock_OrderInfoForUpdate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *OrdersMock_OrderInfoForUpdate_Call) Return(_a0 *model.Order, _a1 error) *OrdersMock_OrderInfoForUpdate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrdersMock_OrderInfoForUpdate_Call) RunAndReturn(run func(context.Context, string) (*model.Order, error)) *OrdersMock_OrderInfoForUpdate_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateOrder provides a mock function with given fields: ctx, order
func (_m *OrdersMock) UpdateOrder(ctx context.Context, order model.Order) error {
	ret := _m.Called(ctx, order)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.Order) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OrdersMock_UpdateOrder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateOrder'
type OrdersMock_UpdateOrder_Call struct {
	*mock.Call
}

// UpdateOrder is a helper method to define mock.On call
//   - ctx context.Context
//   - order model.Order
func (_e *OrdersMock_Expecter) UpdateOrder(ctx interface{}, order interface{}) *OrdersMock_UpdateOrder_Call {
	return &OrdersMock_UpdateOrder_Call{Call: _e.mock.On("UpdateOrder", ctx, order)}
}

func (_c *OrdersMock_UpdateOrder_Call) Run(run func(ctx context.Context, order model.Order)) *OrdersMock_UpdateOrder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.Order))
	})
	return _c
}

func (_c *OrdersMock_UpdateOrder_Call) Return(_a0 error) *OrdersMock_UpdateOrder_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OrdersMock_UpdateOrder_Call) RunAndReturn(run func(context.Context, model.Order) error) *OrdersMock_UpdateOrder_Call {
	_c.Call.Return(run)
	return _c
}

// NewOrdersMock creates a new instance of OrdersMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrdersMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrdersMock {
	mock := &OrdersMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
