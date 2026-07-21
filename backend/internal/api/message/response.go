package message

import (
	"encoding/json"
	"net/http"
)

type ResponseBody struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

func Response(w http.ResponseWriter, status string, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(ResponseBody{
		Status: status,
		Data:   data,
	})
	return err
}
