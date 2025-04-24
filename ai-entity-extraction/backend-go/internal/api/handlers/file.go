package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"backend-go/internal/config"
	"backend-go/internal/services"
	"backend-go/internal/storage"
	"backend-go/pkg/models"

	"github.com/gorilla/mux"
)

// FileHandler gestionează operațiile legate de fișiere
type FileHandler struct {
	cfg       *config.Config
	blobStore *storage.BlobStorage
	sqlStore  *storage.SQLStorage
	iaService *services.EntityExtractor
}

// NewFileHandler creează un nou handler pentru fișiere
func NewFileHandler(cfg *config.Config, blobStore *storage.BlobStorage, sqlStore *storage.SQLStorage, iaService *services.EntityExtractor) *FileHandler {
	return &FileHandler{
		cfg:       cfg,
		blobStore: blobStore,
		sqlStore:  sqlStore,
		iaService: iaService,
	}
}

// UploadFile gestionează încărcarea unui fișier
func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Maxim 10MB per fișier
	const maxUploadSize = 10 * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// Parsează form-ul pentru încărcare
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, fmt.Sprintf("Fișierul este prea mare (max %d MB): %v", maxUploadSize/(1024*1024), err), http.StatusBadRequest)
		return
	}

	// Obține fișierul din formular
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Nu s-a putut obține fișierul: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Verifică tipul fișierului
	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/") && contentType != "application/json" {
		http.Error(w, "Tip de fișier nepermis. Sunt acceptate doar fișiere text și JSON.", http.StatusBadRequest)
		return
	}

	// Generează un nume de fișier unic
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)
	safeName := filepath.Base(filename)

	// Încarcă fișierul în Azure Blob Storage
	blobURL, err := h.blobStore.UploadFile(r.Context(), safeName, file, contentType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Eroare la încărcarea fișierului: %v", err), http.StatusInternalServerError)
		return
	}

	// Citește fișierul pentru a-l procesa cu serviciul de IA
	fileContent, err := h.getFileContent(r.Context(), safeName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Eroare la citirea fișierului: %v", err), http.StatusInternalServerError)
		return
	}

	// Salvează înregistrarea în baza de date
	fileRecord := models.FileRecord{
		FileName:    fileHeader.Filename,
		BlobURL:     blobURL,
		UploadedAt:  time.Now(),
		StatusInfo:  "Uploaded",
		ContentType: contentType,
		FileSize:    fileHeader.Size,
	}

	fileID, err := h.sqlStore.SaveFileRecord(r.Context(), fileRecord)
	if err != nil {
		http.Error(w, fmt.Sprintf("Eroare la salvarea înregistrării: %v", err), http.StatusInternalServerError)
		return
	}

	// Procesează fișierul cu serviciul de IA în fundal
	go h.processFileWithIA(fileID, fileContent)

	// Returnează răspunsul
	response := models.FileUploadResponse{
		FileID:   fileID,
		FileName: fileHeader.Filename,
		BlobURL:  blobURL,
		Message:  "Fișier încărcat cu succes. Procesarea începe în fundal.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// getFileContent obține conținutul unui fișier
func (h *FileHandler) getFileContent(ctx context.Context, filename string) (string, error) {
	// Obține fișierul din Azure Blob Storage
	reader, err := h.blobStore.GetFile(ctx, filename)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	// Citește conținutul fișierului
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// processFileWithIA procesează un fișier cu serviciul de IA
func (h *FileHandler) processFileWithIA(fileID int64, content string) {
	// Procesează fișierul cu serviciul de extragere a entităților
	result, err := h.iaService.ExtractEntities(context.Background(), content, "")
	if err != nil {
		// Actualizează baza de date cu eroarea
		h.sqlStore.UpdateFileProcessResult(context.Background(), fileID, fmt.Sprintf("Eroare la procesare: %v", err))
		return
	}

	// Convertește rezultatul în JSON
	jsonResult, err := json.Marshal(result)
	if err != nil {
		h.sqlStore.UpdateFileProcessResult(context.Background(), fileID, fmt.Sprintf("Eroare la serializarea rezultatului: %v", err))
		return
	}

	// Actualizează baza de date cu rezultatul
	h.sqlStore.UpdateFileProcessResult(context.Background(), fileID, string(jsonResult))
}

// GetFileByID obține un fișier după ID
func (h *FileHandler) GetFileByID(w http.ResponseWriter, r *http.Request) {
	// Obține ID-ul din URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID invalid", http.StatusBadRequest)
		return
	}

	// Obține fișierul din baza de date
	file, err := h.sqlStore.GetFileByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Eroare la obținerea fișierului: %v", err), http.StatusNotFound)
		return
	}

	// Returnează răspunsul
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(file)
}

// ListFiles listează toate fișierele
func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	// Obține parametrii de paginare
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // Implicit
	offset := 0 // Implicit

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err == nil && o >= 0 {
			offset = o
		}
	}

	// Obține fișierele din baza de date
	files, err := h.sqlStore.ListFiles(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Eroare la listarea fișierelor: %v", err), http.StatusInternalServerError)
		return
	}

	// Obține numărul total de fișiere
	count, err := h.sqlStore.CountFiles(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Eroare la numărarea fișierelor: %v", err), http.StatusInternalServerError)
		return
	}

	// Returnează răspunsul
	response := models.FilesListResponse{
		Files: files,
		Count: count,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
