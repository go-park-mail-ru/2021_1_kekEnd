// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_users is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// CheckEmailUnique mocks base method.
func (m *MockUserRepository) CheckEmailUnique(newEmail string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckEmailUnique", newEmail)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckEmailUnique indicates an expected call of CheckEmailUnique.
func (mr *MockUserRepositoryMockRecorder) CheckEmailUnique(newEmail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckEmailUnique", reflect.TypeOf((*MockUserRepository)(nil).CheckEmailUnique), newEmail)
}

// CheckPassword mocks base method.
func (m *MockUserRepository) CheckPassword(password string, user *models.User) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPassword", password, user)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckPassword indicates an expected call of CheckPassword.
func (mr *MockUserRepositoryMockRecorder) CheckPassword(password, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPassword", reflect.TypeOf((*MockUserRepository)(nil).CheckPassword), password, user)
}

// CheckUnsubscribed mocks base method.
func (m *MockUserRepository) CheckUnsubscribed(subscriber, user string) (error, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUnsubscribed", subscriber, user)
	ret0, _ := ret[0].(error)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// CheckUnsubscribed indicates an expected call of CheckUnsubscribed.
func (mr *MockUserRepositoryMockRecorder) CheckUnsubscribed(subscriber, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUnsubscribed", reflect.TypeOf((*MockUserRepository)(nil).CheckUnsubscribed), subscriber, user)
}

// CreateUser mocks base method.
func (m *MockUserRepository) CreateUser(user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepositoryMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), user)
}

// GetModels mocks base method.
func (m *MockUserRepository) GetModels(subs []string, startIndex int) ([]*models.UserNoPassword, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModels", subs, startIndex)
	ret0, _ := ret[0].([]*models.UserNoPassword)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetModels indicates an expected call of GetModels.
func (mr *MockUserRepositoryMockRecorder) GetModels(subs, startIndex interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModels", reflect.TypeOf((*MockUserRepository)(nil).GetModels), subs, startIndex)
}

// GetSubscribers mocks base method.
func (m *MockUserRepository) GetSubscribers(startIndex int, user string) (int, []*models.UserNoPassword, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscribers", startIndex, user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].([]*models.UserNoPassword)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetSubscribers indicates an expected call of GetSubscribers.
func (mr *MockUserRepositoryMockRecorder) GetSubscribers(startIndex, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscribers", reflect.TypeOf((*MockUserRepository)(nil).GetSubscribers), startIndex, user)
}

// GetSubscriptions mocks base method.
func (m *MockUserRepository) GetSubscriptions(startIndex int, user string) (int, []*models.UserNoPassword, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscriptions", startIndex, user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].([]*models.UserNoPassword)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetSubscriptions indicates an expected call of GetSubscriptions.
func (mr *MockUserRepositoryMockRecorder) GetSubscriptions(startIndex, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptions", reflect.TypeOf((*MockUserRepository)(nil).GetSubscriptions), startIndex, user)
}

// GetUserByUsername mocks base method.
func (m *MockUserRepository) GetUserByUsername(username string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", username)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockUserRepositoryMockRecorder) GetUserByUsername(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockUserRepository)(nil).GetUserByUsername), username)
}

// Subscribe mocks base method.
func (m *MockUserRepository) Subscribe(subscriber, user string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", subscriber, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockUserRepositoryMockRecorder) Subscribe(subscriber, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockUserRepository)(nil).Subscribe), subscriber, user)
}

// Unsubscribe mocks base method.
func (m *MockUserRepository) Unsubscribe(subscriber, user string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unsubscribe", subscriber, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unsubscribe indicates an expected call of Unsubscribe.
func (mr *MockUserRepositoryMockRecorder) Unsubscribe(subscriber, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockUserRepository)(nil).Unsubscribe), subscriber, user)
}

// UpdateUser mocks base method.
func (m *MockUserRepository) UpdateUser(user *models.User, change models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", user, change)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserRepositoryMockRecorder) UpdateUser(user, change interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserRepository)(nil).UpdateUser), user, change)
}