package server

import (
	configparser "companies/cmd/internal/configParser"
	"companies/cmd/tests/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewRESTfulServer_NotNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockEventSender := mocks.NewMockEventSender(ctrl)

	httpCfg := configparser.HTTP{
		Addr: "127.0.0.1",
		Port: "8080",
	}

	srv := NewRESTfulServer(httpCfg, mockDB, mockEventSender)

	assert.NotNil(t, srv)
}
