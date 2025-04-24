package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"backend-go/internal/config"
	"backend-go/pkg/models"
	
)

// SQLStorage reprezintă un serviciu pentru interacțiunea cu Azure SQL Database
type SQLStorage struct {
	db  *sql.DB
	cfg *config.Config
}

// NewSQLStorage creează o nouă instanță a serviciului SQLStorage
func NewSQLStorage(cfg *config.Config) (*SQLStorage, error) {
	// Construiește string-ul de conexiune
	connString := fmt.Sprintf("server=%s.database.windows.net;user id=%s;password=%s;database=%s",
		cfg.AzureSQLServerName, cfg.AzureSQLUsername, cfg.AzureSQLPassword, cfg.AzureSQLDBName)

	// Deschide conexiunea la baza de date
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Verifică conexiunea
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// Tabelul FileProcessingHistory există deja în baza de date
	// Nu mai e nevoie să-l inițializăm

	return &SQLStorage{
		db:  db,
		cfg: cfg,
	}, nil
}

// SaveFileRecord salvează o înregistrare pentru un fișier în baza de date
func (s *SQLStorage) SaveFileRecord(ctx context.Context, file models.FileRecord) (int64, error) {
	// SQL pentru inserare conform structurii existente
	query := `
	INSERT INTO [dbo].[FileProcessingHistory] (FileName, BlobUrl, UploadTimestamp, ProcessingResult, Status)
	VALUES (@fileName, @blobURL, @uploadTimestamp, @processingResult, @status);
	SELECT SCOPE_IDENTITY();
	`

	// Execută inserarea
	var id int64
	err := s.db.QueryRowContext(
		ctx,
		query,
		sql.Named("fileName", file.FileName),
		sql.Named("blobURL", file.BlobURL),
		sql.Named("uploadTimestamp", time.Now()),
		sql.Named("processingResult", nil), // Inițial este null
		sql.Named("status", "Uploaded"),    // Status inițial
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("error saving file record: %v", err)
	}

	return id, nil
}

// UpdateFileProcessResult actualizează rezultatul procesării pentru un fișier
func (s *SQLStorage) UpdateFileProcessResult(ctx context.Context, fileID int64, result string) error {
	// SQL pentru actualizare conform structurii existente
	query := `
	UPDATE [dbo].[FileProcessingHistory]
	SET ProcessingResult = @processingResult, Status = @status
	WHERE Id = @fileID
	`

	// Execută actualizarea
	_, err := s.db.ExecContext(
		ctx,
		query,
		sql.Named("processingResult", result),
		sql.Named("status", "Processed"),
		sql.Named("fileID", fileID),
	)

	if err != nil {
		return fmt.Errorf("error updating file process result: %v", err)
	}

	return nil
}

// GetFileByID obține un fișier după ID conform structurii existente
func (s *SQLStorage) GetFileByID(ctx context.Context, fileID int64) (*models.FileRecord, error) {
	// SQL pentru selectare
	query := `
	SELECT Id, FileName, BlobUrl, UploadTimestamp, ProcessingResult, Status
	FROM [dbo].[FileProcessingHistory]
	WHERE Id = @fileID
	`

	// Execută interogarea
	row := s.db.QueryRowContext(ctx, query, sql.Named("fileID", fileID))

	var file models.FileRecord
	var processingResult sql.NullString
	var status string

	err := row.Scan(
		&file.ID,
		&file.FileName,
		&file.BlobURL,
		&file.UploadedAt,
		&processingResult,
		&status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("error getting file: %v", err)
	}

	if processingResult.Valid {
		file.ProcessResult = processingResult.String
	}
	
	// Setăm StatusInfo pentru a include statusul în răspunsul nostru
	file.StatusInfo = status

	return &file, nil
}

// ListFiles returnează toate fișierele din baza de date
func (s *SQLStorage) ListFiles(ctx context.Context, limit, offset int) ([]models.FileRecord, error) {
	// SQL pentru listare conform structurii existente
	query := `
	SELECT Id, FileName, BlobUrl, UploadTimestamp, ProcessingResult, Status
	FROM [dbo].[FileProcessingHistory]
	ORDER BY UploadTimestamp DESC
	OFFSET @offset ROWS
	FETCH NEXT @limit ROWS ONLY
	`

	// Execută interogarea
	rows, err := s.db.QueryContext(
		ctx,
		query,
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)
	if err != nil {
		return nil, fmt.Errorf("error listing files: %v", err)
	}
	defer rows.Close()

	var files []models.FileRecord
	for rows.Next() {
		var file models.FileRecord
		var processingResult sql.NullString
		var status string

		err := rows.Scan(
			&file.ID,
			&file.FileName,
			&file.BlobURL,
			&file.UploadedAt,
			&processingResult,
			&status,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning file row: %v", err)
		}

		if processingResult.Valid {
			file.ProcessResult = processingResult.String
		}
		
		file.StatusInfo = status

		files = append(files, file)
	}

	return files, nil
}

// CountFiles returnează numărul total de fișiere din baza de date
func (s *SQLStorage) CountFiles(ctx context.Context) (int, error) {
	// SQL pentru numărare
	query := `SELECT COUNT(*) FROM [dbo].[FileProcessingHistory]`

	// Execută interogarea
	var count int
	err := s.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting files: %v", err)
	}

	return count, nil
}

// Close închide conexiunea la baza de date
func (s *SQLStorage) Close() error {
	return s.db.Close()
}