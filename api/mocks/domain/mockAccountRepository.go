// Code generated by MockGen. DO NOT EDIT.
// Source: red/domain (interfaces: AccountRepository)

// Package domain is a generated GoMock package.
package domain

import (
	domain "red/domain"
	errs "red/errs"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAccountRepository is a mock of AccountRepository interface.
type MockAccountRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAccountRepositoryMockRecorder
}

// MockAccountRepositoryMockRecorder is the mock recorder for MockAccountRepository.
type MockAccountRepositoryMockRecorder struct {
	mock *MockAccountRepository
}

// NewMockAccountRepository creates a new mock instance.
func NewMockAccountRepository(ctrl *gomock.Controller) *MockAccountRepository {
	mock := &MockAccountRepository{ctrl: ctrl}
	mock.recorder = &MockAccountRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountRepository) EXPECT() *MockAccountRepositoryMockRecorder {
	return m.recorder
}

// ByCustomerID mocks base method.
func (m *MockAccountRepository) ByCustomerID(arg0 uint) ([]domain.Account, *errs.AppError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByCustomerID", arg0)
	ret0, _ := ret[0].([]domain.Account)
	ret1, _ := ret[1].(*errs.AppError)
	return ret0, ret1
}

// ByCustomerID indicates an expected call of ByCustomerID.
func (mr *MockAccountRepositoryMockRecorder) ByCustomerID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByCustomerID", reflect.TypeOf((*MockAccountRepository)(nil).ByCustomerID), arg0)
}

// ByID mocks base method.
func (m *MockAccountRepository) ByID(arg0 uint) (*domain.Account, *errs.AppError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByID", arg0)
	ret0, _ := ret[0].(*domain.Account)
	ret1, _ := ret[1].(*errs.AppError)
	return ret0, ret1
}

// ByID indicates an expected call of ByID.
func (mr *MockAccountRepositoryMockRecorder) ByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByID", reflect.TypeOf((*MockAccountRepository)(nil).ByID), arg0)
}

// Delete mocks base method.
func (m *MockAccountRepository) Delete(arg0 string) *errs.AppError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(*errs.AppError)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAccountRepositoryMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAccountRepository)(nil).Delete), arg0)
}

// Save mocks base method.
func (m *MockAccountRepository) Save(arg0 domain.Account) (*domain.Account, *errs.AppError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0)
	ret0, _ := ret[0].(*domain.Account)
	ret1, _ := ret[1].(*errs.AppError)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockAccountRepositoryMockRecorder) Save(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockAccountRepository)(nil).Save), arg0)
}

// SaveTransaction mocks base method.
func (m *MockAccountRepository) SaveTransaction(arg0 domain.Transaction) (*domain.Transaction, *errs.AppError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTransaction", arg0)
	ret0, _ := ret[0].(*domain.Transaction)
	ret1, _ := ret[1].(*errs.AppError)
	return ret0, ret1
}

// SaveTransaction indicates an expected call of SaveTransaction.
func (mr *MockAccountRepositoryMockRecorder) SaveTransaction(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTransaction", reflect.TypeOf((*MockAccountRepository)(nil).SaveTransaction), arg0)
}

// TransactionsByID mocks base method.
func (m *MockAccountRepository) TransactionsByID(arg0 uint) ([]domain.Transaction, *errs.AppError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransactionsByID", arg0)
	ret0, _ := ret[0].([]domain.Transaction)
	ret1, _ := ret[1].(*errs.AppError)
	return ret0, ret1
}

// TransactionsByID indicates an expected call of TransactionsByID.
func (mr *MockAccountRepositoryMockRecorder) TransactionsByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionsByID", reflect.TypeOf((*MockAccountRepository)(nil).TransactionsByID), arg0)
}

// Update mocks base method.
func (m *MockAccountRepository) Update(arg0 domain.Account) (*domain.Account, *errs.AppError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(*domain.Account)
	ret1, _ := ret[1].(*errs.AppError)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockAccountRepositoryMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAccountRepository)(nil).Update), arg0)
}
