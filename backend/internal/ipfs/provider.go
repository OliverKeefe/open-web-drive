package ipfs

import (
	"context"
	"io"

	"github.com/ipfs/go-cid"
)

// IPFSProvider interface for IPFS operations.
type IPFSProvider interface {
	Authorize(ctx context.Context, principal string) (string, error)
	Upload(ctx context.Context, rdr io.Reader) (cid.Cid, error)
	UnPin(ctx context.Context, c cid.Cid) (bool, error)
	CheckStatus(ctx context.Context, c cid.Cid) (bool, error)
}
