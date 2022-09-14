// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/sdk/keycloak/keycloak.go

// Package mock_keycloak is a generated GoMock package.
package mock_keycloak

import (
	context "context"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	keycloak "github.com/indrasaputra/arjuna/pkg/sdk/keycloak"
)

// MockDoer is a mock of Doer interface.
type MockDoer struct {
	ctrl     *gomock.Controller
	recorder *MockDoerMockRecorder
}

// MockDoerMockRecorder is the mock recorder for MockDoer.
type MockDoerMockRecorder struct {
	mock *MockDoer
}

// NewMockDoer creates a new mock instance.
func NewMockDoer(ctrl *gomock.Controller) *MockDoer {
	mock := &MockDoer{ctrl: ctrl}
	mock.recorder = &MockDoerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDoer) EXPECT() *MockDoerMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockDoer) Do(arg0 *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockDoerMockRecorder) Do(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockDoer)(nil).Do), arg0)
}

// MockKeycloak is a mock of Keycloak interface.
type MockKeycloak struct {
	ctrl     *gomock.Controller
	recorder *MockKeycloakMockRecorder
}

// MockKeycloakMockRecorder is the mock recorder for MockKeycloak.
type MockKeycloakMockRecorder struct {
	mock *MockKeycloak
}

// NewMockKeycloak creates a new mock instance.
func NewMockKeycloak(ctrl *gomock.Controller) *MockKeycloak {
	mock := &MockKeycloak{ctrl: ctrl}
	mock.recorder = &MockKeycloakMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeycloak) EXPECT() *MockKeycloakMockRecorder {
	return m.recorder
}

// CreateClient mocks base method.
func (m *MockKeycloak) CreateClient(ctx context.Context, token, realm string, client *keycloak.ClientRepresentation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateClient", ctx, token, realm, client)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateClient indicates an expected call of CreateClient.
func (mr *MockKeycloakMockRecorder) CreateClient(ctx, token, realm, client interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateClient", reflect.TypeOf((*MockKeycloak)(nil).CreateClient), ctx, token, realm, client)
}

// CreateRealm mocks base method.
func (m *MockKeycloak) CreateRealm(ctx context.Context, token string, realm *keycloak.RealmRepresentation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRealm", ctx, token, realm)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRealm indicates an expected call of CreateRealm.
func (mr *MockKeycloakMockRecorder) CreateRealm(ctx, token, realm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRealm", reflect.TypeOf((*MockKeycloak)(nil).CreateRealm), ctx, token, realm)
}

// CreateUser mocks base method.
func (m *MockKeycloak) CreateUser(ctx context.Context, token, realm string, user *keycloak.UserRepresentation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, token, realm, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockKeycloakMockRecorder) CreateUser(ctx, token, realm, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockKeycloak)(nil).CreateUser), ctx, token, realm, user)
}

// GetUserByEmail mocks base method.
func (m *MockKeycloak) GetUserByEmail(ctx context.Context, token, realm, email string) (*keycloak.UserRepresentation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, token, realm, email)
	ret0, _ := ret[0].(*keycloak.UserRepresentation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockKeycloakMockRecorder) GetUserByEmail(ctx, token, realm, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockKeycloak)(nil).GetUserByEmail), ctx, token, realm, email)
}

// LoginAdmin mocks base method.
func (m *MockKeycloak) LoginAdmin(ctx context.Context, username, password string) (*keycloak.JWT, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginAdmin", ctx, username, password)
	ret0, _ := ret[0].(*keycloak.JWT)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginAdmin indicates an expected call of LoginAdmin.
func (mr *MockKeycloakMockRecorder) LoginAdmin(ctx, username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginAdmin", reflect.TypeOf((*MockKeycloak)(nil).LoginAdmin), ctx, username, password)
}
