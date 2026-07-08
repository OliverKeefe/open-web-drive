package files

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestHandler_ScheduleDelete(t *testing.T) {
	h := svc.Handle
	payload := DeleteRequest{ID: uuid.New()}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal("invalid test json payload")
	}

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/files/schedule-deletion",
		bytes.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")

	recorder := httptest.NewRecorder()

	h(recorder, req)
	if recorder.Code != 204 {
		t.Fatalf("expected 204, got %v", recorder.Code)
	}
}
