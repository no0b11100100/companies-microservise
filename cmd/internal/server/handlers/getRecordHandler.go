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

type getRecordDB interface {
	GetRecord(uuid.UUID) (database.CompanyInfo, error)
}

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
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(record)
	}
}

// func NewGetRecordHandler(db getRecordDB, eventSender eventsender.EventSender) Handler {
// 	h := getRecordHandler{}
// 	base := initBaseHandler(db, h.handler, eventSender)
// 	h.baseHandler = base

// 	return &h
// }

// func (c *getRecordHandler) handler(w http.ResponseWriter, r *http.Request) {
// 	log.Println(consts.ApplicationPrefix, "Request path: ", r.URL.Path)
// 	log.Println(consts.ApplicationPrefix, "Path param id: ", chi.URLParam(r, "id"))

// 	id1 := chi.URLParam(r, "id")
// 	fmt.Println("ID:", id1)
// 	id, err := uuid.Parse(id1)
// 	if err != nil {
// 		log.Println(consts.ApplicationPrefix, "getRecordHandler::handler error:", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	record, err := c.db.GetRecord(id)

// 	if err != nil {
// 		log.Println(consts.ApplicationPrefix, "getRecordHandler::handler error:", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(record)
// }
