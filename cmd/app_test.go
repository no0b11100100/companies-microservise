package main

import (
	"companies/cmd/tests/mocks"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCloser struct {
	mock.Mock
}

func (m *MockCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestApp_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockREST := mocks.NewMockRESTServer(ctrl)
	mockREST.EXPECT().Serve().Times(1)

	a := &app{
		restServer: mockREST,
	}

	a.Run()
}

func TestApp_Close_AllOK(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockDB := &MockCloser{}
	mockREST := mocks.NewMockRESTServer(ctrl)
	mockEvent := &MockCloser{}

	mockDB.On("Close").Return(nil)
	mockEvent.On("Close").Return(nil)
	mockREST.EXPECT().Shutdown().Return(nil)

	a := &app{
		db:          mockDB,
		restServer:  mockREST,
		eventSender: mockEvent,
	}

	err := a.Close()
	assert.NoError(t, err)

	mockDB.AssertCalled(t, "Close")
	mockEvent.AssertCalled(t, "Close")
}

func TestApp_Close_WithErrors(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockDB := &MockCloser{}
	mockREST := mocks.NewMockRESTServer(ctrl)
	mockEvent := &MockCloser{}

	mockDB.On("Close").Return(errors.New("db error"))
	mockEvent.On("Close").Return(errors.New("sender error"))
	mockREST.EXPECT().Shutdown().Return(errors.New("shutdown err"))

	a := &app{
		db:          mockDB,
		restServer:  mockREST,
		eventSender: mockEvent,
	}

	err := a.Close()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db err")
	assert.Contains(t, err.Error(), "shutdown err")
	assert.Contains(t, err.Error(), "event err")
}
