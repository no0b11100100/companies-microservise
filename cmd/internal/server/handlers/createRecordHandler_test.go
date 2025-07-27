package handlers

import (
	"bytes"
	"companies/cmd/internal/database"
	"companies/cmd/tests/mocks"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func makeValidCompany() database.CompanyInfo {
	name := "Test Company"
	desc := "Some description"
	emp := 100
	isReg := true
	typ := 1
	id := uuid.New()

	return database.CompanyInfo{
		ID:             &id,
		Name:           &name,
		Description:    &desc,
		EmployeesCount: &emp,
		IsRegistered:   &isReg,
		Type:           &typ,
	}
}

func TestNewCreateRecordHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockDB := mocks.NewMockcreateRecordDB(ctrl)
	mockSender := mocks.NewMockEventSender(ctrl)

	company := makeValidCompany()

	mockDB.EXPECT().IsRecordExists(*company.Name).Return(false)
	mockDB.EXPECT().CreateRecord(company).Return(*company.ID, nil)
	mockSender.EXPECT().PublishEvent("data-changed", gomock.Any()).Times(1)

	body, _ := json.Marshal(company)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/companies", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler := NewCreateRecordHandler(mockDB, mockSender)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %v", rr.Code)
	}
}

func TestNewCreateRecordHandler_InvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockDB := mocks.NewMockcreateRecordDB(ctrl)
	mockSender := mocks.NewMockEventSender(ctrl)

	invalid := database.CompanyInfo{}

	mockSender.EXPECT().PublishEvent("data-changed", gomock.Any()).Times(1)

	body, _ := json.Marshal(invalid)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/companies", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler := NewCreateRecordHandler(mockDB, mockSender)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %v", rr.Code)
	}
}

func TestNewCreateRecordHandler_RecordExists(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockDB := mocks.NewMockcreateRecordDB(ctrl)
	mockSender := mocks.NewMockEventSender(ctrl)

	company := makeValidCompany()

	mockDB.EXPECT().IsRecordExists(*company.Name).Return(true)
	mockSender.EXPECT().PublishEvent("data-changed", gomock.Any()).Times(1)

	body, _ := json.Marshal(company)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/companies", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler := NewCreateRecordHandler(mockDB, mockSender)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %v", rr.Code)
	}
}

func TestNewCreateRecordHandler_DBCreateError(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockDB := mocks.NewMockcreateRecordDB(ctrl)
	mockSender := mocks.NewMockEventSender(ctrl)

	company := makeValidCompany()

	mockDB.EXPECT().IsRecordExists(*company.Name).Return(false)
	mockDB.EXPECT().CreateRecord(company).Return(uuid.Nil, context.DeadlineExceeded)
	mockSender.EXPECT().PublishEvent("data-changed", gomock.Any()).Times(1)

	body, _ := json.Marshal(company)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/companies", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler := NewCreateRecordHandler(mockDB, mockSender)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %v", rr.Code)
	}
}
