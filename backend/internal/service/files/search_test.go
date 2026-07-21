package files

// TODO: Monolithic svc was refactored into SearchService.
// This test needs to be rewritten to test SearchService directly.

/*
import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_FindAll(t *testing.T) {
	var cur = MetadataCursor{
		ModifiedAt: time.Now(),
		ID:         uuid.New(),
	}

	var req = GetAllMetadataRequest{
		UserID: uuid.New(),
		Cursor: &cur,
		Limit:  20,
	}

	md, err := svc.FindAllMetadata(context.Background(), req)
	...
}
*/
