package entities

import (
	"encoding/json"
)

type request struct {
	// BUG: if id is number - in response it would be string :c
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
}

func NewRequest(data []byte) (request, error) {
	var req request
	err := json.Unmarshal(data, &req)
	return req, err
}
