package handlers

import (
	"companies/cmd/internal/consts"
	"companies/cmd/internal/database"
	eventsender "companies/cmd/internal/eventSender"
	"encoding/json"
	"log"
	"net/http"

	"companies/cmd/internal/structs"

	"github.com/google/uuid"
)

const (
	kMaxNameLenght        = 15
	kMaxDescriptionLenght = 3000
)

//go:generate mockgen -source=createRecordHandler.go -destination=../../../tests/mocks/mock_create_record.go -package=mocks
type createRecordDB interface {
	CreateRecord(database.CompanyInfo) (uuid.UUID, error)
	IsRecordExists(string) bool
}

func IsValidInfo(data database.CompanyInfo) bool {
	if data.Name == nil {
		return false
	}

	if len(*data.Name) > kMaxNameLenght {
		return false
	}

	if data.Description != nil && len(*data.Description) > kMaxDescriptionLenght {
		return false
	}

	if data.EmployeesCount == nil {
		return false
	}

	if data.IsRegistered == nil {
		return false
	}

	if data.Type == nil {
		return false
	}

	return true
}

// @Summary      Create a new company record
// @Description  Creates a new company with the provided information
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        company  body      database.CompanyInfo  true  "Company to create"
// @Success      201      {object}  map[string]string      "Created. Returns the new company ID"
// @Failure      400      {string}  string                 "Bad request – invalid input or error"
// @Failure      409      {string}  string                 "Conflict – record already exists"
// @Router       /api/v1/companies [post]
func NewCreateRecordHandler(db createRecordDB, eventSender eventsender.EventSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(consts.ApplicationPrefix, "createRecordHandler::handler", r.Body)
		var record database.CompanyInfo
		json.NewDecoder(r.Body).Decode(&record)

		if !IsValidInfo(record) {
			log.Println(consts.ApplicationPrefix, "createRecordHandler::handler invalid data")
			w.WriteHeader(http.StatusBadRequest)
			eventSender.PublishEvent("data-changed", structs.Event{
				URL:           r.URL.Path,
				Type:          structs.Created,
				Status:        structs.Failed,
				ErrorMesssage: "invalid data provided",
			})
			return
		}

		if db.IsRecordExists(*record.Name) {
			log.Println(consts.ApplicationPrefix, "createRecordHandler::handler record alredy exist")
			w.WriteHeader(http.StatusConflict)
			eventSender.PublishEvent("data-changed", structs.Event{
				URL:           r.URL.Path,
				Type:          structs.Created,
				Status:        structs.Failed,
				ErrorMesssage: "record alredy exist",
			})
			return
		}

		id, err := db.CreateRecord(record)
		if err != nil {
			eventSender.PublishEvent("data-changed", structs.Event{
				URL:           r.URL.Path,
				Type:          structs.Created,
				Status:        structs.Failed,
				ErrorMesssage: err.Error(),
			})
			log.Println(consts.ApplicationPrefix, "createRecordHandler::handler error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		eventSender.PublishEvent("data-changed", structs.Event{
			URL:    r.URL.Path,
			Type:   structs.Created,
			Status: structs.Success,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		resp := map[string]string{"companyId": id.String()}
		json.NewEncoder(w).Encode(resp)
	}
}
