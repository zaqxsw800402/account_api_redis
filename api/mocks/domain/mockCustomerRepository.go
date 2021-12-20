// Code generated by MockGen. DO NOT EDIT.
// Source: red/domain (interfaces: CustomerRepository)

// Package domain is a generated GoMock package.
package domain

import (
	domain "red/domain"
	errs "red/errs"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCustomerRepository is a mock of CustomerRepository interface.
type MockCustomerRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCustomerRepositoryMockRecorder
}

// MockCustomerRepositoryMockRecorder is the mock recorder for MockCustomerRepository.
type MockCustomerRepositoryMockRecorder struct {
	mock *MockCustomerRepository
}

// NewMockCustomerRepository creates a new mock instance.
func NewMockCustomerRepository(ctrl *gomock.Controller) *MockCustomerRepository {
	mock := &MockCustomerRepository{ctrl: ctrl}
	mock.recorder = &MockCustomerRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCustomerRepository) EXPECT() *MockCustomerRepositoryMockRecorder {
	return m.recorder
}

// ByID mocks base method.
func (m *MockCustomerRepository) ByID(arg0 string) (*domain.Customer, *errs.AppError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByID", arg0)
	ret0, _ := ret[0].(*domain.Customer)
	ret1, _ := ret[1].(*errs.AppError)
	return ret0, ret1
}

// ByID indicates an expected call of ByID.
func (mr *MockCustomerRepositoryMockRecorder) ByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByID", reflect.TypeOf((*MockCustomerRepository)(nil).ByID), arg0)
}

// Delete mocks base method.
func (m *MockCustomerRepository) Delete(arg0 string) *errs.AppError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(*errs.AppError)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCustomerRepositoryMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCustomerRepository)(nil).Delete), arg0)
}

// FindAll mocks base method.
func (m *MockCustomerRepository) FindAll(arg0 int) ([]domain.Customer, *errs.AppError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", arg0)
	ret0, _ := ret[0].([]domain.Customer)
	ret1, _ := ret[1].(*errs.AppError)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockCustomerRepositoryMockRecorder) FindAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockCustomerRepository)(nil).FindAll), arg0)
}

// Save mocks base method.
func (m *MockCustomerRepository) Save(arg0 domain.Customer) (*domain.Customer, *errs.AppError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0)
	ret0, _ := ret[0].(*domain.Customer)
	ret1, _ := ret[1].(*errs.AppError)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockCustomerRepositoryMockRecorder) Save(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockCustomerRepository)(nil).Save), arg0)
}

// Update mocks base method.
func (m *MockCustomerRepository) Update(arg0 domain.Customer) (*domain.Customer, *errs.AppError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(*domain.Customer)
	ret1, _ := ret[1].(*errs.AppError)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockCustomerRepositoryMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCustomerRepository)(nil).Update), arg0)
}
