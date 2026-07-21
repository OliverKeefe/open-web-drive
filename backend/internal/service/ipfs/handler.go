package ipfs

import (
	"backend/internal/ipfs"
	"backend/internal/service/files"
	"context"
)

type service interface {
}

type Handler struct {
	svc  service
	ipfs ipfs.IPFSProvider
}

func (h *Handler) UploadToIPFS(ctx context.Context, file []files.File) {

}
