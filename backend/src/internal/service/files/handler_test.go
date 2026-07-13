package files

// TODO: Handler type was refactored into individual service types
// (SearchService, DeleteService, etc.). These tests need to be
// rewritten to test the new service types directly.

/*
import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockService struct{}

var (
	mockSvc = &mockService{}
)

func (m *mockService) Upload(ctx context.Context, r *http.Request) ([]MetaData, error) {
	return []MetaData{}, nil
}

func (m *mockService) FindMetadata(ctx context.Context, request FindMetadataRequest) ([]MetaData, error) {
	return []MetaData{}, nil
}

func (m *mockService) Delete(ctx context.Context, request DeleteRequest) error {
	return nil
}

func (m *mockService) MoveToRubbish(ctx context.Context, request DeleteRequest) error {
	return nil
}

func (m *mockService) FindAllMetadata(ctx context.Context, request GetAllMetadataRequest) ([]MetaDataResponse, error) {
	return []MetaDataResponse{}, nil
}

func TestHandler_GetAll(t *testing.T) {
	h := &Handler{svc: mockSvc}
	...
}

func TestHandler_Delete(t *testing.T) {
	h := Handler{svc: mockSvc}
	...
}
*/
