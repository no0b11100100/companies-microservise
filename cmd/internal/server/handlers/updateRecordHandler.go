package handlers

import (
	"companies/cmd/internal/consts"
	"companies/cmd/internal/database"
	eventsender "companies/cmd/internal/eventSender"
	"companies/cmd/internal/structs"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

//go:generate mockgen -source=updateRecordHandler.go -destination=../../../tests/mocks/mock_update_record.go -package=mocks
type updateRecordDB interface {
	UpdateRecord(database.CompanyInfo, uuid.UUID) error
}

// @Summary      Update an existing company
// @Description  Updates company information by UUID
// @Tags         Companies
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id      path      string                true  "Company UUID"
// @Param        company body      database.CompanyInfo  true  "Updated company data"
// @Success      202     {string}  string                "Accepted – update in progress"
// @Failure      400     {string}  string                "Bad request – invalid UUID or body"
// @Router       /api/v1/companies/{id} [patch]
func NewUpdateRecordHandler(db updateRecordDB, eventSender eventsender.EventSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(consts.ApplicationPrefix, "updateRecordHandler::handler", r.Body)

		log.Println(consts.ApplicationPrefix, "Request path: ", r.URL.Path)
		log.Println(consts.ApplicationPrefix, "Path param id: ", chi.URLParam(r, "id"))

		uuidStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(uuidStr)

		if err != nil {
			log.Println(consts.ApplicationPrefix, "updateRecordHandler::handler error:", err)
			eventSender.PublishEvent("data-changed", structs.Event{
				URL:           r.URL.Path,
				Type:          structs.Updated,
				Status:        structs.Failed,
				ErrorMesssage: err.Error(),
			})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data := database.CompanyInfo{}

		json.NewDecoder(r.Body).Decode(&data)

		err = db.UpdateRecord(data, id)

		if err != nil {
			eventSender.PublishEvent("data-changed", structs.Event{
				URL:           r.URL.Path,
				Type:          structs.Updated,
				Status:        structs.Failed,
				ErrorMesssage: err.Error(),
			})
			log.Println(consts.ApplicationPrefix, "updateRecordHandler::handler error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bytes, _ := json.Marshal(data)
		eventSender.PublishEvent("data-changed", structs.Event{
			URL:    r.URL.Path,
			Type:   structs.Updated,
			Status: structs.Success,
			Data:   bytes,
		})

		w.WriteHeader(http.StatusAccepted)
	}
}
