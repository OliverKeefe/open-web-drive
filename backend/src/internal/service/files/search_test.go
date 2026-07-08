package files

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
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(md) != reflect.TypeOf([]MetaData{}) {
		t.Fatalf("expected slice of Metadata, got: %T", md)
	}
}
