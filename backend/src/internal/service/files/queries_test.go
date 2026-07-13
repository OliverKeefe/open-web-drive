package files

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
)

func TestQuery_FindAllMetadata_NoCursor(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	repo := NewFileRepository(mock)

	userID := uuid.New()
	limit := 2
	now := time.Now()

	rows := pgxmock.NewRows([]string{
		"id", "file_name", "path", "size", "file_type",
		"modified_at", "uploaded_at", "owner_id",
		"checksum", "version",
	}).
		AddRow(
			uuid.New(), "a.txt", "dir1/dir2/a.txt", uint64(10),
			".txt", now, now, userID, []byte("deadbeef"), now,
		).
		AddRow(
			uuid.New(), "b.png", "/b.png", uint64(64),
			".png", now, now, userID, []byte("deadbeef"), now,
		)

	mock.ExpectQuery(`SELECT id, file_name, path, size, file_type, modified_at,`).
		WithArgs(userID, limit).
		WillReturnRows(rows)

	res, err := repo.FindAllMetadata(ctx, GetAllMetadataRequest{
		UserID: userID,
		Limit:  limit,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(res))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestQuery_FindAllMetadata_Cursor(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	repo := NewFileRepository(mock)

	userID := uuid.New()
	cursorID := uuid.New()
	now := time.Now()
	limit := 2
	cur := MetadataCursor{
		ModifiedAt: now,
		ID:         cursorID,
	}

	rows := pgxmock.NewRows([]string{
		"id", "file_name", "path", "size", "file_type",
		"modified_at", "uploaded_at", "owner_id",
		"checksum", "version",
	}).
		AddRow(uuid.New(), "a.txt", "dir1/dir2/a.txt", uint64(10), ".txt",
			now, now, userID, []byte("deadbeef"), now,
		).
		AddRow(uuid.New(), "c.java", "src/c.java", uint64(64), ".java",
			now, now, userID, []byte("deadbeef"), now,
		)

	mock.ExpectQuery(`SELECT id, file_name, path, size, file_type, modified_at,`).
		WithArgs(userID, cur.ModifiedAt, cur.ID, limit).
		WillReturnRows(rows)

	res, err := repo.FindAllMetadata(ctx, GetAllMetadataRequest{
		UserID: userID,
		Cursor: &cur,
		Limit:  limit,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(res))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
