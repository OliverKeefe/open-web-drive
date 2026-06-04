package files

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Metadata struct {
	File   FileMetadata
	Access AccessMetadata
	Ipfs   IPFSMetadata
}

// Metadata Model - need ContentCID string for IPFS
type MetaData struct {
	ID         uuid.UUID   `json:"uuid"`
	FileName   string      `json:"file_name"`
	Path       string      `json:"path"`
	Size       uint64      `json:"size"`
	FileType   string      `json:"file_type"`
	ModifiedAt time.Time   `json:"modified_at"`
	UploadedAt time.Time   `json:"created_at"`
	Owner      uuid.UUID   `json:"owner_id"`
	AccessTo   []uuid.UUID `json:"access_to"`
	Group      []uuid.UUID `json:"group_id"`
	CheckSum   []byte      `json:"checksum"`
	Version    time.Time   `json:"version"`
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

type AccessMetadata struct {
	OwnerID     uuid.UUID   `json:"owner_id"`
	SharedWith  []uuid.UUID `json:"shared_with"`
	GroupAccess []uuid.UUID `json:"group_access"`
}

type FileMetadata struct {
	ID           uuid.UUID         `json:"uuid"`
	FileName     string            `json:"file_name"`
	Path         string            `json:"path"`
	RelativePath string            `json:"relative_path"`
	Size         uint64            `json:"size"`
	FileType     string            `json:"file_type"`
	ModifiedAt   time.Time         `json:"modified_at"`
	UploadedAt   time.Time         `json:"uploaded_at"`
	CreatedAt    time.Time         `json:"created_at"`
	Version      uuid.UUID         `json:"version"`
	Hash         [32]byte          `json:"hash"`
	Permissions  []FilePermissions `json:"permissions"`
}

type FilePermissions struct {
	//FileID       uuid.UUID `json:"file_id"`
	FileID      uuid.UUID `json:"permissions_id"`
	GranteeType string    `json:"grantee_type"`
	GranteeID   uuid.UUID `json:"grantee_id"`
	AccessLevel uuid.UUID `json:"access_level"`
}

func (m *MetaData) ToResponse() MetaDataResponse {
	return MetaDataResponse{
		ID:         m.ID,
		FileName:   m.FileName,
		Path:       m.Path,
		Size:       m.Size,
		FileType:   m.FileType,
		ModifiedAt: m.ModifiedAt,
		UploadedAt: m.UploadedAt,
		Owner:      m.Owner,
		AccessTo:   m.AccessTo,
		Group:      m.Group,
		CheckSum:   m.CheckSum,
		Version:    m.Version,
	}
}

type MetaDataResponse struct {
	ID         uuid.UUID   `json:"uuid"`
	FileName   string      `json:"file_name"`
	Path       string      `json:"path"`
	Size       uint64      `json:"size"`
	FileType   string      `json:"file_type"`
	ModifiedAt time.Time   `json:"modified_at"`
	UploadedAt time.Time   `json:"created_at"`
	Owner      uuid.UUID   `json:"owner_id"`
	AccessTo   []uuid.UUID `json:"access_to"`
	Group      []uuid.UUID `json:"group_id"`
	CheckSum   []byte      `json:"checksum"`
	Version    time.Time   `json:"version"`
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
	Size       uint64      `json:"size,omitempty"`
	FileType   string      `json:"file_type,omitempty"`
	ModifiedAt time.Time   `json:"modified_at,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at,omitempty"`
	Owner      uuid.UUID   `json:"owner_id"`
	AccessTo   []uuid.UUID `json:"access_to,omitempty"`
	Group      []uuid.UUID `json:"group_id,omitempty"`
	CheckSum   []byte      `json:"checksum,omitempty"`
	Version    time.Time   `json:"version,omitempty"`
}

func (req *FindMetadataRequest) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	return nil
}

func (req *FindMetadataRequest) ToModel() MetaData {
	return MetaData{
		ID:         req.ID,
		FileName:   req.FileName,
		Path:       req.Path,
		Size:       req.Size,
		FileType:   req.FileType,
		ModifiedAt: req.ModifiedAt,
		UploadedAt: req.UploadedAt,
		Owner:      req.Owner,
		AccessTo:   req.AccessTo,
		Group:      req.Group,
		CheckSum:   req.CheckSum,
		Version:    req.Version,
	}
}

type DeleteRequest struct {
	ID uuid.UUID `json:"id"`
	//OwnerID uuid.UUID `json:"owner_id"`
}

func (req *DeleteRequest) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	return nil
}
