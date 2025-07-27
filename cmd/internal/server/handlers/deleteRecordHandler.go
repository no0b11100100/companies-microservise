package handlers

import (
	"companies/cmd/internal/consts"
	eventsender "companies/cmd/internal/eventSender"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type deleteRecordDB interface {
	DeleteRecord(uuid.UUID) error
}

// @Summary      Delete a company
// @Description  Deletes a company record by its UUID
// @Tags         Companies
// @Security     BearerAuth
// @Param        id   path      string  true  "Company UUID"
// @Success      200  {string}  string  "Successfully deleted"
// @Failure      400  {string}  string  "Invalid UUID or deletion failed"
// @Router       /api/v1/companies/{id} [delete]
func NewDeleteRecordHandler(db deleteRecordDB, eventSender eventsender.EventSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(consts.ApplicationPrefix, "deleteRecordHandler::handler", r.Body)

		log.Println(consts.ApplicationPrefix, "Request path: ", r.URL.Path)
		log.Println(consts.ApplicationPrefix, "Path param id: ", chi.URLParam(r, "id"))

		uuidStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(uuidStr)

		if err != nil {
			eventSender.PublishEvent("data-changed", eventsender.Event{
				URL:           r.URL.Path,
				Type:          eventsender.Deleted,
				Status:        eventsender.Failed,
				ErrorMesssage: err.Error(),
			})
			log.Println(consts.ApplicationPrefix, "deleteRecordHandler::handler error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := db.DeleteRecord(id); err != nil {
			log.Println(consts.ApplicationPrefix, "deleteRecordHandler::handler error:", err)
			eventSender.PublishEvent("data-changed", eventsender.Event{
				URL:           r.URL.Path,
				Type:          eventsender.Deleted,
				Status:        eventsender.Failed,
				ErrorMesssage: err.Error(),
			})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		eventSender.PublishEvent("data-changed", eventsender.Event{
			URL:    r.URL.Path,
			Type:   eventsender.Deleted,
			Status: eventsender.Success,
		})

		w.WriteHeader(http.StatusOK)
	}
}
