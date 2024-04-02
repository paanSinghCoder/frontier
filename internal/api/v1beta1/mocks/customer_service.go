// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	customer "github.com/raystack/frontier/billing/customer"
	mock "github.com/stretchr/testify/mock"
)

// CustomerService is an autogenerated mock type for the CustomerService type
type CustomerService struct {
	mock.Mock
}

type CustomerService_Expecter struct {
	mock *mock.Mock
}

func (_m *CustomerService) EXPECT() *CustomerService_Expecter {
	return &CustomerService_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, _a1
func (_m *CustomerService) Create(ctx context.Context, _a1 customer.Customer) (customer.Customer, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 customer.Customer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, customer.Customer) (customer.Customer, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, customer.Customer) customer.Customer); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(customer.Customer)
	}

	if rf, ok := ret.Get(1).(func(context.Context, customer.Customer) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CustomerService_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type CustomerService_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 customer.Customer
func (_e *CustomerService_Expecter) Create(ctx interface{}, _a1 interface{}) *CustomerService_Create_Call {
	return &CustomerService_Create_Call{Call: _e.mock.On("Create", ctx, _a1)}
}

func (_c *CustomerService_Create_Call) Run(run func(ctx context.Context, _a1 customer.Customer)) *CustomerService_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(customer.Customer))
	})
	return _c
}

func (_c *CustomerService_Create_Call) Return(_a0 customer.Customer, _a1 error) *CustomerService_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CustomerService_Create_Call) RunAndReturn(run func(context.Context, customer.Customer) (customer.Customer, error)) *CustomerService_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, id
func (_m *CustomerService) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CustomerService_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type CustomerService_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *CustomerService_Expecter) Delete(ctx interface{}, id interface{}) *CustomerService_Delete_Call {
	return &CustomerService_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *CustomerService_Delete_Call) Run(run func(ctx context.Context, id string)) *CustomerService_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *CustomerService_Delete_Call) Return(_a0 error) *CustomerService_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CustomerService_Delete_Call) RunAndReturn(run func(context.Context, string) error) *CustomerService_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *CustomerService) GetByID(ctx context.Context, id string) (customer.Customer, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 customer.Customer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (customer.Customer, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) customer.Customer); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(customer.Customer)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CustomerService_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type CustomerService_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *CustomerService_Expecter) GetByID(ctx interface{}, id interface{}) *CustomerService_GetByID_Call {
	return &CustomerService_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *CustomerService_GetByID_Call) Run(run func(ctx context.Context, id string)) *CustomerService_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *CustomerService_GetByID_Call) Return(_a0 customer.Customer, _a1 error) *CustomerService_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CustomerService_GetByID_Call) RunAndReturn(run func(context.Context, string) (customer.Customer, error)) *CustomerService_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, filter
func (_m *CustomerService) List(ctx context.Context, filter customer.Filter) ([]customer.Customer, error) {
	ret := _m.Called(ctx, filter)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []customer.Customer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, customer.Filter) ([]customer.Customer, error)); ok {
		return rf(ctx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, customer.Filter) []customer.Customer); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]customer.Customer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, customer.Filter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CustomerService_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type CustomerService_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - filter customer.Filter
func (_e *CustomerService_Expecter) List(ctx interface{}, filter interface{}) *CustomerService_List_Call {
	return &CustomerService_List_Call{Call: _e.mock.On("List", ctx, filter)}
}

func (_c *CustomerService_List_Call) Run(run func(ctx context.Context, filter customer.Filter)) *CustomerService_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(customer.Filter))
	})
	return _c
}

func (_c *CustomerService_List_Call) Return(_a0 []customer.Customer, _a1 error) *CustomerService_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CustomerService_List_Call) RunAndReturn(run func(context.Context, customer.Filter) ([]customer.Customer, error)) *CustomerService_List_Call {
	_c.Call.Return(run)
	return _c
}

// ListPaymentMethods provides a mock function with given fields: ctx, id
func (_m *CustomerService) ListPaymentMethods(ctx context.Context, id string) ([]customer.PaymentMethod, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for ListPaymentMethods")
	}

	var r0 []customer.PaymentMethod
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]customer.PaymentMethod, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []customer.PaymentMethod); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]customer.PaymentMethod)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CustomerService_ListPaymentMethods_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListPaymentMethods'
type CustomerService_ListPaymentMethods_Call struct {
	*mock.Call
}

// ListPaymentMethods is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *CustomerService_Expecter) ListPaymentMethods(ctx interface{}, id interface{}) *CustomerService_ListPaymentMethods_Call {
	return &CustomerService_ListPaymentMethods_Call{Call: _e.mock.On("ListPaymentMethods", ctx, id)}
}

func (_c *CustomerService_ListPaymentMethods_Call) Run(run func(ctx context.Context, id string)) *CustomerService_ListPaymentMethods_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *CustomerService_ListPaymentMethods_Call) Return(_a0 []customer.PaymentMethod, _a1 error) *CustomerService_ListPaymentMethods_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CustomerService_ListPaymentMethods_Call) RunAndReturn(run func(context.Context, string) ([]customer.PaymentMethod, error)) *CustomerService_ListPaymentMethods_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, _a1
func (_m *CustomerService) Update(ctx context.Context, _a1 customer.Customer) (customer.Customer, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 customer.Customer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, customer.Customer) (customer.Customer, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, customer.Customer) customer.Customer); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(customer.Customer)
	}

	if rf, ok := ret.Get(1).(func(context.Context, customer.Customer) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CustomerService_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type CustomerService_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 customer.Customer
func (_e *CustomerService_Expecter) Update(ctx interface{}, _a1 interface{}) *CustomerService_Update_Call {
	return &CustomerService_Update_Call{Call: _e.mock.On("Update", ctx, _a1)}
}

func (_c *CustomerService_Update_Call) Run(run func(ctx context.Context, _a1 customer.Customer)) *CustomerService_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(customer.Customer))
	})
	return _c
}

func (_c *CustomerService_Update_Call) Return(_a0 customer.Customer, _a1 error) *CustomerService_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CustomerService_Update_Call) RunAndReturn(run func(context.Context, customer.Customer) (customer.Customer, error)) *CustomerService_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewCustomerService creates a new instance of CustomerService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCustomerService(t interface {
	mock.TestingT
	Cleanup(func())
}) *CustomerService {
	mock := &CustomerService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}