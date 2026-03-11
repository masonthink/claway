package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/clawbeach/server/internal/model"
	"github.com/jackc/pgx/v5"
)

// GetDocumentByTaskID retrieves the document associated with a task.
func (s *Store) GetDocumentByTaskID(ctx context.Context, taskID int64) (*model.Document, error) {
	var d model.Document
	err := s.db.QueryRow(ctx,
		`SELECT id, task_id, content, current_version, created_at, updated_at
		 FROM documents WHERE task_id = $1`, taskID,
	).Scan(&d.ID, &d.TaskID, &d.Content, &d.CurrentVersion, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get document by task_id: %w", err)
	}
	return &d, nil
}

// CreateDocument creates a new empty document for a task.
func (s *Store) CreateDocument(ctx context.Context, taskID int64) (*model.Document, error) {
	var d model.Document
	err := s.db.QueryRow(ctx,
		`INSERT INTO documents (task_id)
		 VALUES ($1)
		 RETURNING id, task_id, content, current_version, created_at, updated_at`, taskID,
	).Scan(&d.ID, &d.TaskID, &d.Content, &d.CurrentVersion, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create document: %w", err)
	}
	return &d, nil
}

// UpdateDocumentContent updates the content and version of a document.
func (s *Store) UpdateDocumentContent(ctx context.Context, docID int64, content string, version int) error {
	tag, err := s.db.Exec(ctx,
		`UPDATE documents SET content = $1, current_version = $2, updated_at = NOW()
		 WHERE id = $3`, content, version, docID,
	)
	if err != nil {
		return fmt.Errorf("update document content: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// CreateDocumentVersion inserts a new version snapshot for a document.
func (s *Store) CreateDocumentVersion(ctx context.Context, docID int64, version int, content, diff string, createdBy int64) (*model.DocumentVersion, error) {
	var dv model.DocumentVersion
	err := s.db.QueryRow(ctx,
		`INSERT INTO document_versions (document_id, version, content, diff_from_previous, created_by)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, document_id, version, content, diff_from_previous, created_at, created_by`,
		docID, version, content, diff, createdBy,
	).Scan(&dv.ID, &dv.DocumentID, &dv.Version, &dv.Content, &dv.DiffFromPrevious, &dv.CreatedAt, &dv.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("create document version: %w", err)
	}
	return &dv, nil
}

// ListDocumentVersions returns all versions of a document ordered by version number.
func (s *Store) ListDocumentVersions(ctx context.Context, docID int64) ([]*model.DocumentVersion, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, document_id, version, content, diff_from_previous, created_at, created_by
		 FROM document_versions WHERE document_id = $1
		 ORDER BY version ASC`, docID,
	)
	if err != nil {
		return nil, fmt.Errorf("list document versions: %w", err)
	}
	defer rows.Close()

	var versions []*model.DocumentVersion
	for rows.Next() {
		var dv model.DocumentVersion
		if err := rows.Scan(&dv.ID, &dv.DocumentID, &dv.Version, &dv.Content, &dv.DiffFromPrevious, &dv.CreatedAt, &dv.CreatedBy); err != nil {
			return nil, fmt.Errorf("list document versions scan: %w", err)
		}
		versions = append(versions, &dv)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list document versions rows: %w", err)
	}
	return versions, nil
}

// GetDocumentVersion retrieves a specific version of a document.
func (s *Store) GetDocumentVersion(ctx context.Context, docID int64, version int) (*model.DocumentVersion, error) {
	var dv model.DocumentVersion
	err := s.db.QueryRow(ctx,
		`SELECT id, document_id, version, content, diff_from_previous, created_at, created_by
		 FROM document_versions WHERE document_id = $1 AND version = $2`, docID, version,
	).Scan(&dv.ID, &dv.DocumentID, &dv.Version, &dv.Content, &dv.DiffFromPrevious, &dv.CreatedAt, &dv.CreatedBy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get document version: %w", err)
	}
	return &dv, nil
}
