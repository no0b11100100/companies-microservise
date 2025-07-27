package handlers

import (
	"companies/cmd/internal/database"
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

func ptrString(s string) *string {
	return &s
}

func TestGetRecordHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockgetRecordDB(ctrl)
	handler := NewGetRecordHandler(mockDB)

	id := uuid.New()
	record := database.CompanyInfo{
		ID:   &id,
		Name: ptrString("Test Company"),
	}

	mockDB.EXPECT().GetRecord(id).Return(record, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/companies/"+id.String(), nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var got database.CompanyInfo
	err := json.NewDecoder(rr.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, *record.ID, *got.ID)
	assert.Equal(t, *record.Name, *got.Name)
}

func TestGetRecordHandler_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockgetRecordDB(ctrl)
	handler := NewGetRecordHandler(mockDB)

	invalidID := "invalid-uuid"
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", invalidID)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/companies/"+invalidID, nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetRecordHandler_RecordNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockgetRecordDB(ctrl)
	handler := NewGetRecordHandler(mockDB)

	id := uuid.New()

	mockDB.EXPECT().GetRecord(id).Return(database.CompanyInfo{}, errors.New("not found"))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/companies/"+id.String(), nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
