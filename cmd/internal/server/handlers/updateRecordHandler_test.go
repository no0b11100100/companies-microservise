package handlers

import (
	"bytes"
	"companies/cmd/internal/database"
	"companies/cmd/internal/structs"
	"companies/cmd/tests/mocks"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateRecordHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockupdateRecordDB(ctrl)
	mockEventSender := mocks.NewMockEventSender(ctrl)

	handler := NewUpdateRecordHandler(mockDB, mockEventSender)

	id := uuid.New()
	company := database.CompanyInfo{
		Name: ptrString("Updated Company"),
	}

	// Ожидаем вызов UpdateRecord без ошибки
	mockDB.EXPECT().UpdateRecord(company, id).Return(nil)

	// Ожидаем вызов PublishEvent с успешным статусом
	mockEventSender.EXPECT().PublishEvent("data-changed", gomock.Any()).DoAndReturn(
		func(topic string, event structs.Event) error {
			assert.Equal(t, structs.Updated, event.Type)
			assert.Equal(t, structs.Success, event.Status)
			return nil
		})

	bodyBytes, _ := json.Marshal(company)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/companies/"+id.String(), bytes.NewReader(bodyBytes))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusAccepted, rr.Code)
}

func TestUpdateRecordHandler_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockupdateRecordDB(ctrl)
	mockEventSender := mocks.NewMockEventSender(ctrl)

	handler := NewUpdateRecordHandler(mockDB, mockEventSender)

	invalidID := "not-a-uuid"
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", invalidID)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/companies/"+invalidID, nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Ожидаем вызов PublishEvent с ошибкой
	mockEventSender.EXPECT().PublishEvent("data-changed", gomock.Any()).Return(nil)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateRecordHandler_UpdateRecordError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockupdateRecordDB(ctrl)
	mockEventSender := mocks.NewMockEventSender(ctrl)

	handler := NewUpdateRecordHandler(mockDB, mockEventSender)

	id := uuid.New()
	company := database.CompanyInfo{
		Name: ptrString("Company With Error"),
	}

	mockDB.EXPECT().UpdateRecord(company, id).Return(errors.New("update failed"))

	mockEventSender.EXPECT().PublishEvent("data-changed", gomock.Any()).Return(nil)

	bodyBytes, _ := json.Marshal(company)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/companies/"+id.String(), bytes.NewReader(bodyBytes))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
