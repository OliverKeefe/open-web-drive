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
		"id", "file_id", "file_name", "path", "relative_path", "size", "file_type",
		"owner_id", "version", "hash", "created_at", "modified_at", "uploaded_at",
	}).
		AddRow(
			uuid.New(), uuid.New(), "a.txt", "dir1/dir2/a.txt", "/a.txt", int64(10),
			".txt", userID, 1, "deadbeef", now, now, now,
		).
		AddRow(
			uuid.New(), uuid.New(), "b.png", "/b.png", "/b.png", int64(64),
			".png", userID, 1, "deadbeef", now, now, now,
		)

	mock.ExpectQuery(`SELECT id, file_id, file_name, path, relative_path, size, file_type`).
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
		"id", "file_id", "file_name", "path", "relative_path", "size", "file_type",
		"owner_id", "version", "hash", "created_at", "modified_at", "uploaded_at",
	}).
		AddRow(uuid.New(), uuid.New(), "a.txt", "dir1/dir2/a.txt", "/a.txt", int64(10), ".txt",
			userID, 1, "deadbeef", now, now, now,
		).
		AddRow(uuid.New(), uuid.New(), "c.java", "src/c.java", "/c.java", int64(64), ".java",
			userID, 2, "deadbeef", now, now, now,
		)

	mock.ExpectQuery(`SELECT id, file_id, file_name, path, relative_path, size, file_type`).
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

var metadataColumns = []string{
	"id", "file_id", "file_name", "path", "relative_path", "size", "file_type",
	"owner_id", "version", "hash", "created_at", "modified_at", "uploaded_at",
}

func newMetadataRow() *pgxmock.Rows {
	return pgxmock.NewRows(metadataColumns)
}

func TestQuery_FindMetadata_NoFilters(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewFileRepository(mock)
	ctx := context.Background()
	now := time.Now()
	uid := uuid.New()

	mock.ExpectQuery(`SELECT id, file_id, file_name, path, relative_path, size, file_type`).
		WillReturnRows(newMetadataRow().
			AddRow(uid, uuid.New(), "doc.pdf", "docs/doc.pdf", "/doc.pdf", int64(2048), ".pdf",
				uid, 1, "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", now, now, now))

	res, err := repo.FindMetadata(ctx, FileMetadata{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 row, got %d", len(res))
	}
	if res[0].FileName != "doc.pdf" {
		t.Fatalf("expected file_name doc.pdf, got %s", res[0].FileName)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestQuery_FindMetadata_SingleFilter(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewFileRepository(mock)
	ctx := context.Background()
	now := time.Now()
	uid := uuid.New()

	mock.ExpectQuery(`SELECT id, file_id, file_name, path, relative_path, size, file_type`).
		WithArgs(uid).
		WillReturnRows(newMetadataRow().
			AddRow(uuid.New(), uuid.New(), "img.png", "/img.png", "/img.png", int64(4096), ".png",
				uid, 1, "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", now, now, now))

	res, err := repo.FindMetadata(ctx, FileMetadata{OwnerID: uid})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 row, got %d", len(res))
	}
	if res[0].FileType != ".png" {
		t.Fatalf("expected file_type .png, got %s", res[0].FileType)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestQuery_FindMetadata_MultipleFilters(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewFileRepository(mock)
	ctx := context.Background()
	now := time.Now()
	uid := uuid.New()
	fileID := uuid.New()

	mock.ExpectQuery(`SELECT id, file_id, file_name, path, relative_path, size, file_type`).
		WithArgs(fileID, ".pdf", uid, 2).
		WillReturnRows(newMetadataRow().
			AddRow(uuid.New(), fileID, "report.pdf", "docs/report.pdf", "/report.pdf", int64(8192), ".pdf",
				uid, 2, "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", now, now, now))

	res, err := repo.FindMetadata(ctx, FileMetadata{
		OwnerID:  uid,
		FileID:   fileID,
		FileType: ".pdf",
		Version:  2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 row, got %d", len(res))
	}
	if res[0].Version != 2 {
		t.Fatalf("expected version 2, got %d", res[0].Version)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestQuery_FindMetadata_EmptyResult(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewFileRepository(mock)
	ctx := context.Background()
	uid := uuid.New()

	mock.ExpectQuery(`SELECT id, file_id, file_name, path, relative_path, size, file_type`).
		WithArgs(uid).
		WillReturnRows(newMetadataRow())

	res, err := repo.FindMetadata(ctx, FileMetadata{OwnerID: uid})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(res))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
