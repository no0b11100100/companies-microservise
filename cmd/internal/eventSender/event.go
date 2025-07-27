package eventsender

import "encoding/json"

const (
	Created int = iota
	Updated
	Deleted
)

const (
	Success = iota
	Failed
)

type Event struct {
	Type          int             `json:"type"`
	Status        int             `json:"status"`
	Data          json.RawMessage `json:"data,omitempty"`
	ErrorMesssage string          `json:"errorMessage,omitempty"`
}
