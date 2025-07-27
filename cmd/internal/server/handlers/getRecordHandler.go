package handlers

import (
	"companies/cmd/internal/consts"
	"companies/cmd/internal/database"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

//go:generate mockgen -source=getRecordHandler.go -destination=../../../tests/mocks/mock_get_record.go -package=mocks
type getRecordDB interface {
	GetRecord(uuid.UUID) (database.CompanyInfo, error)
}

// @Summary      Get a company by ID
// @Description  Retrieves company information using a UUID
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        id   path      string                true  "Company UUID"
// @Success      200  {object}  database.CompanyInfo  "Company found"
// @Failure      400  {string}  string                "Invalid UUID"
// @Failure      404  {string}  string                "Company not found"
// @Router       /api/v1/companies/{id} [get]
func NewGetRecordHandler(db getRecordDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(consts.ApplicationPrefix, "getRecordHandler::handler")

		log.Println(consts.ApplicationPrefix, "Request path: ", r.URL.Path)
		log.Println(consts.ApplicationPrefix, "Path param id: ", chi.URLParam(r, "id"))

		uuidStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(uuidStr)
		if err != nil {
			log.Println(consts.ApplicationPrefix, "getRecordHandler::handler error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		record, err := db.GetRecord(id)

		if err != nil {
			log.Println(consts.ApplicationPrefix, "getRecordHandler::handler error:", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(record)
	}
}
