package files

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// FileMetadata is the canonical model matching the file_metadata table.
type FileMetadata struct {
	ID           uuid.UUID `json:"id"`
	FileID       uuid.UUID `json:"file_id"`
	Version      int       `json:"version"`
	OwnerID      uuid.UUID `json:"owner_id"`
	FileName     string    `json:"file_name"`
	Path         string    `json:"path"`
	RelativePath string    `json:"relative_path"`
	Size         int64     `json:"size"`
	FileType     string    `json:"file_type"`
	Hash         string    `json:"hash"`
	CreatedAt    time.Time `json:"created_at"`
	ModifiedAt   time.Time `json:"modified_at"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

// FileMetadataResponse is the API response model.
type FileMetadataResponse = FileMetadata

// FilePermissions matches the file_permissions table.
type FilePermissions struct {
	ID            uuid.UUID `json:"id"`
	FileID        uuid.UUID `json:"file_id"`
	GranteeType   string    `json:"grantee_type"`
	GranteeID     uuid.UUID `json:"grantee_id"`
	AccessLevelID uuid.UUID `json:"access_level_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type Visibility int

const (
	Public Visibility = iota
	Private
)

type IPFSMetadata struct {
	CID            string     `json:"cid"`
	Space          string     `json:"space"`
	Visibility     Visibility `json:"visibility"`
	DID            string     `json:"did"`
	Shards         []string   `json:"shards"`
	Uri            string     `json:"uri"`
	HTTPGatewayURL string     `json:"http_gateway_url"`
}

type GetAllMetadataRequest struct {
	UserID uuid.UUID       `json:"user_id"`
	Cursor *MetadataCursor `json:"cursor"`
	Limit  int             `json:"limit"`
}

func (req *GetAllMetadataRequest) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("invalid json request: %w", err)
	}
	return nil
}

type MetadataCursor struct {
	ModifiedAt time.Time `json:"modified_at"`
	ID         uuid.UUID `json:"id"`
}

type FindMetadataRequest struct {
	ID         uuid.UUID   `json:"file_id"`
	FileName   string      `json:"file_name,omitempty"`
	Path       string      `json:"path,omitempty"`
	Size       int64       `json:"size,omitempty"`
	FileType   string      `json:"file_type,omitempty"`
	ModifiedAt time.Time   `json:"modified_at,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at,omitempty"`
	Owner      uuid.UUID   `json:"owner_id"`
	AccessTo   []uuid.UUID `json:"access_to,omitempty"`
	Group      []uuid.UUID `json:"group_id,omitempty"`
	Hash       string      `json:"hash,omitempty"`
	Version    int         `json:"version,omitempty"`
}

func (req *FindMetadataRequest) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	return nil
}

func (req *FindMetadataRequest) ToModel() FileMetadata {
	return FileMetadata{
		ID:         req.ID,
		FileName:   req.FileName,
		Path:       req.Path,
		Size:       req.Size,
		FileType:   req.FileType,
		ModifiedAt: req.ModifiedAt,
		UploadedAt: req.UploadedAt,
		OwnerID:    req.Owner,
		Hash:       req.Hash,
		Version:    req.Version,
	}
}

type DeleteRequest struct {
	ID uuid.UUID `json:"id"`
}

func (req *DeleteRequest) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	return nil
}

type DownloadRequest struct {
	ID uuid.UUID `json:"id"`
}

func (req *DownloadRequest) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return fmt.Errorf("invalid json request: %w", err)
	}
	return nil
}
