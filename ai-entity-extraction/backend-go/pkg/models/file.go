package models

import (
	"time"
)

// FileRecord reprezintă o înregistrare din baza de date pentru un fișier încărcat
type FileRecord struct {
	ID            int64     `json:"id"`
	FileName      string    `json:"fileName"`
	BlobURL       string    `json:"blobUrl"`
	UploadedAt    time.Time `json:"uploadedAt"`
	ProcessResult string    `json:"processingResult,omitempty"`
	StatusInfo    string    `json:"status"`
	ContentType   string    `json:"contentType,omitempty"` // Pentru compatibilitate
	FileSize      int64     `json:"fileSize,omitempty"`    // Pentru compatibilitate
}

// ProcessResult reprezintă rezultatul procesării prin serviciul de IA
type ProcessResult struct {
	Entities []Entity `json:"entities"`
	Success  bool     `json:"success"`
	Message  string   `json:"message,omitempty"`
}

// Entity reprezintă o entitate extrasă din text
type Entity struct {
	Name       string   `json:"name"`
	Category   string   `json:"category"`
	Confidence float64  `json:"confidence"`
	Offset     int      `json:"offset"`
	Length     int      `json:"length"`
	SubType    string   `json:"subType,omitempty"`
	Matches    []string `json:"matches,omitempty"`
}

// FileUploadResponse reprezintă răspunsul la încărcarea unui fișier
type FileUploadResponse struct {
	FileID   int64  `json:"fileId"`
	FileName string `json:"fileName"`
	BlobURL  string `json:"blobUrl"`
	Message  string `json:"message"`
}

// FilesListResponse reprezintă răspunsul la listarea fișierelor
type FilesListResponse struct {
	Files []FileRecord `json:"files"`
	Count int          `json:"count"`
}
