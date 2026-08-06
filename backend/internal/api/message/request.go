package message

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Bind[T any](r *http.Request) (T, error) {
	var request T
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, fmt.Errorf("invalid json request: %w", err)
	}
	return request, nil
}
