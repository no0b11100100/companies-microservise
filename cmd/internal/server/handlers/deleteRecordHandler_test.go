package handlers

import (
	"bytes"
	// "companies/cmd/internal/eventsender"

	"companies/cmd/internal/structs"
	"companies/cmd/tests/mocks"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func newDeleteTestRequest(method, url string, id string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(nil))
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", id)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	return req.WithContext(ctx)
}

func TestDeleteRecordHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockdeleteRecordDB(ctrl)
	mockSender := mocks.NewMockEventSender(ctrl)

	testID := uuid.New()

	mockDB.EXPECT().DeleteRecord(testID).Return(nil)
	mockSender.EXPECT().PublishEvent("data-changed", gomock.Any()).Times(1)

	req := newDeleteTestRequest(http.MethodDelete, "/api/v1/companies/"+testID.String(), testID.String())
	rr := httptest.NewRecorder()

	handler := NewDeleteRecordHandler(mockDB, mockSender)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteRecordHandler_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockdeleteRecordDB(ctrl)
	mockSender := mocks.NewMockEventSender(ctrl)

	mockSender.EXPECT().PublishEvent("data-changed", gomock.AssignableToTypeOf(structs.Event{
		Type:   structs.Deleted,
		Status: structs.Failed,
	})).Times(1)

	req := newDeleteTestRequest(http.MethodDelete, "/api/v1/companies/invalid-uuid", "invalid-uuid")
	rr := httptest.NewRecorder()

	handler := NewDeleteRecordHandler(mockDB, mockSender)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteRecordHandler_DeleteFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockdeleteRecordDB(ctrl)
	mockSender := mocks.NewMockEventSender(ctrl)

	testID := uuid.New()
	expectedErr := errors.New("delete failed")

	mockDB.EXPECT().DeleteRecord(testID).Return(expectedErr)
	mockSender.EXPECT().PublishEvent("data-changed", gomock.AssignableToTypeOf(structs.Event{
		Type:   structs.Deleted,
		Status: structs.Failed,
	})).Times(1)

	req := newDeleteTestRequest(http.MethodDelete, "/api/v1/companies/"+testID.String(), testID.String())
	rr := httptest.NewRecorder()

	handler := NewDeleteRecordHandler(mockDB, mockSender)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
