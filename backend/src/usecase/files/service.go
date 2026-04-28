package files

import (
	"backend/src/internal/auth"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type repository interface {
	SaveMetaData(ctx context.Context, meta MetaData) (MetaData, error) //TODO: Change order of params
	SaveFileData(basePath string, rdr io.Reader, filename string) error
	SaveToS3(ctx context.Context, basePath string, rdr io.Reader, filename string) error
	FindMetadata(ctx context.Context, model MetaData) ([]MetaData, error)
	DeleteMetadata(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error
	DeleteFileData(ctx context.Context, fileID uuid.UUID, ownerID uuid.UUID) error
	FindAllMetadata(ctx context.Context, req GetAllMetadataRequest) ([]MetaData, error)
	MarkForDeletion(ctx context.Context, id uuid.UUID, id2 uuid.UUID) error
}

type Service struct {
	repo repository
}

// NewService constructor for new Service (UploadService).
func NewService(repo repository) *Service {
	return &Service{
		repo: repo,
	}
}

// multipartMetadata is a Data Transfer Object for multipart form for file uploads
// from the web gui frontend.
// When the web frontend sends a multipart form, metadata is stored
// as raw json.
type multipartMetadata struct {
	Path             string `json:"path"`
	RelativePath     string `json:"relativePath"`
	LastModified     int64  `json:"lastModified"`
	LastModifiedDate string `json:"lastModifiedDate"`
	Size             uint64 `json:"size"`
	FileType         string `json:"fileType"`

	ID       string `json:"id"`
	OwnerID  string `json:"ownerId"`
	CheckSum []byte `json:"checkSum"`
}

func (svc *Service) Upload(ctx context.Context, r *http.Request) ([]MetaData, error) {
	mr, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}
	debug := false
	metadataByID := make(map[string]MetaData)

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		name := part.FormName()

		switch {
		// Handle Metadata.
		case strings.HasPrefix(name, "metadata-"):
			idStr := strings.TrimPrefix(name, "metadata-")

			// decode + build metadata
			var decodedRequest multipartMetadata
			if err := json.NewDecoder(part).Decode(&decodedRequest); err != nil {
				return nil, err
			}

			userId, ok := auth.UserIDFromCtx(r.Context())
			if !ok {
				return nil, errors.New("could not get userID from context")
			}

			ownerID, err := uuid.Parse(userId)
			if err != nil {
				return nil, err
			}

			metadataByID[idStr] = MetaData{
				ID:         uuid.MustParse(idStr),
				FileName:   decodedRequest.RelativePath,
				Path:       decodedRequest.Path,
				Size:       decodedRequest.Size,
				ModifiedAt: time.UnixMilli(decodedRequest.LastModified),
				UploadedAt: time.Now(),
				Owner:      ownerID,
				Version:    time.Now(),
			}

		// Handle Part containing file's binary data.
		case strings.HasPrefix(name, "file-"):
			idStr := strings.TrimPrefix(name, "file-")
			// File has to be saved here, if you try to pass this to another temp location
			// in memory then the data will be unusable.

			data, err := io.ReadAll(part)
			if err != nil {
				return nil, err
			}

			hash := sha256.Sum256(data)

			if debug == false {
				err = svc.repo.SaveToS3(
					ctx,
					"",
					bytes.NewReader(data),
					part.FileName(),
				)
			}
			if debug == true {

				hash := sha256.New()

				if err := svc.repo.SaveFileData(
					"/home/oliver/Development/gestalt/gestalt/backend/tempfiles",
					io.TeeReader(part, hash),
					part.FileName(),
				); err != nil {
					return nil, err
				}
			}

			md := metadataByID[idStr]
			md.CheckSum = hash[:]
			metadataByID[idStr] = md
		}
	}
	// Persist file metadata
	var newMetadata []MetaData
	for _, md := range metadataByID {
		newMd, err := svc.repo.SaveMetaData(ctx, md)
		if err != nil {
			return nil, err
		}
		newMetadata = append(newMetadata, newMd)
	}

	return newMetadata, nil
}

func (svc *Service) FindAllMetadata(ctx context.Context, request GetAllMetadataRequest) ([]MetaDataResponse, error) {
	var (
		repo     = svc.repo
		files    []MetaData
		response []MetaDataResponse
	)

	files, err := repo.FindAllMetadata(ctx, request)
	if err != nil {
		log.Printf("unable to get all files for user: %s, %v ", request.UserID, err)
	}

	for _, file := range files {
		file := file.ToResponse()
		if err != nil {
			log.Printf("unable to map file metadata: %v, to dto: %v", file, err)
		}
		response = append(response, file)
	}

	return response, nil
}

func (svc *Service) FindMetadata(ctx context.Context, request FindMetadataRequest) ([]MetaData, error) {
	var (
		repo  = svc.repo
		files []MetaData
	)

	model := request.ToModel()
	files, err := repo.FindMetadata(ctx, model)
	if err != nil {
		log.Printf("unable to get file metadata: %v", err)
		return files, err
	}

	return files, nil
}

func (svc *Service) Delete(ctx context.Context, request DeleteRequest) error {
	userID, ok := auth.UserIDFromCtx(ctx)
	if !ok {
		return errors.New("unable to get userID from context")
	}

	ownerID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("unable to parse userID string to uuid")
	}

	err = svc.repo.DeleteMetadata(ctx, request.ID, ownerID)
	if err != nil {
		log.Printf("could not delete file metadata, %v", err)
		return err
	}

	return nil
}

func (svc *Service) MoveToRubbish(ctx context.Context, request DeleteRequest) error {
	userID, ok := auth.UserIDFromCtx(ctx)
	if !ok {
		return errors.New("unable to get userID from context")
	}

	ownerID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("unable to parse userID string to uuid")
	}

	err = svc.repo.MarkForDeletion(ctx, request.ID, ownerID)
	if err != nil {
		log.Printf("unable to move file or metadata to rubbish bin, %v", err)
		return err
	}

	return nil
}
