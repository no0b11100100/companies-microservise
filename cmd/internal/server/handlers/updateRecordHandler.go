package handlers

import (
	"companies/cmd/internal/consts"
	"companies/cmd/internal/database"
	eventsender "companies/cmd/internal/eventSender"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type updateRecordDB interface {
	UpdateRecord(database.CompanyInfo, uuid.UUID) error
}

// @Summary Update an existing company
// @Description Updates the company by ID. Requires JWT authentication.
// @Tags companies
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param company body database.CompanyInfo true "Updated company data"
// @Success 200 {object} database.CompanyInfo
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/companies/{id} [patch]
func NewUpdateRecordHandler(db updateRecordDB, eventSender eventsender.EventSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(consts.ApplicationPrefix, "updateRecordHandler::handler", r.Body)

		log.Println(consts.ApplicationPrefix, "Request path: ", r.URL.Path)
		log.Println(consts.ApplicationPrefix, "Path param id: ", chi.URLParam(r, "id"))

		uuidStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(uuidStr)

		if err != nil {
			log.Println(consts.ApplicationPrefix, "updateRecordHandler::handler error:", err)
			eventSender.PublishEvent("data-changed", eventsender.Event{
				Type:          eventsender.Updated,
				Status:        eventsender.Failed,
				ErrorMesssage: err.Error(),
			})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data := database.CompanyInfo{}

		json.NewDecoder(r.Body).Decode(&data)

		err = db.UpdateRecord(data, id)

		if err != nil {
			eventSender.PublishEvent("data-changed", eventsender.Event{
				Type:          eventsender.Updated,
				Status:        eventsender.Failed,
				ErrorMesssage: err.Error(),
			})
			log.Println(consts.ApplicationPrefix, "updateRecordHandler::handler error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bytes, _ := json.Marshal(data)
		eventSender.PublishEvent("data-changed", eventsender.Event{
			Type:   eventsender.Updated,
			Status: eventsender.Success,
			Data:   bytes,
		})

		w.WriteHeader(http.StatusAccepted)
	}
}
