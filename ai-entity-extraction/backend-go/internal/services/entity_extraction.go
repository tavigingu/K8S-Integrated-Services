package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"backend-go/internal/config"
	"backend-go/pkg/models"
)

// EntityExtractor reprezintă un serviciu pentru extragerea entităților folosind Azure Cognitive Services
type EntityExtractor struct {
	endpoint string
	apiKey   string
}

// NewEntityExtractor creează o nouă instanță a serviciului EntityExtractor
func NewEntityExtractor(cfg *config.Config) *EntityExtractor {
	return &EntityExtractor{
		endpoint: cfg.AzureCognitiveServicesEndpoint,
		apiKey:   cfg.AzureCognitiveServicesKey,
	}
}

// Request reprezintă cererea către API-ul de extragere a entităților
type Request struct {
	Documents []Document `json:"documents"`
}

// Document reprezintă un document pentru procesare
type Document struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	Language string `json:"language"`
}

// Response reprezintă răspunsul de la API-ul de extragere a entităților
type Response struct {
	Documents []DocumentResponse `json:"documents"`
	Errors    []ErrorResponse    `json:"errors"`
}

// DocumentResponse reprezintă răspunsul pentru un document
type DocumentResponse struct {
	ID       string           `json:"id"`
	Entities []EntityResponse `json:"entities"`
}

// EntityResponse reprezintă o entitate extrasă
type EntityResponse struct {
	Name       string          `json:"text"`
	Category   string          `json:"category"`
	Confidence float64         `json:"confidenceScore"`
	Offset     int             `json:"offset"`
	Length     int             `json:"length"`
	SubType    string          `json:"subCategory,omitempty"`
	Matches    []MatchResponse `json:"matches,omitempty"`
}

// MatchResponse reprezintă o potrivire pentru o entitate
type MatchResponse struct {
	Text       string  `json:"text"`
	Offset     int     `json:"offset"`
	Length     int     `json:"length"`
	Confidence float64 `json:"confidenceScore"`
}

// ErrorResponse reprezintă o eroare în răspunsul API-ului
type ErrorResponse struct {
	ID      string `json:"id"`
	Message string `json:"error"`
}

// ExtractEntities extrage entitățile dintr-un text
func (e *EntityExtractor) ExtractEntities(ctx context.Context, text string, language string) (*models.ProcessResult, error) {
	// Construim cererea
	if language == "" {
		language = "en"
	}

	req := Request{
		Documents: []Document{
			{
				ID:       "1",
				Text:     text,
				Language: language,
			},
		},
	}

	// Convertim cererea în JSON
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Pregătim cererea HTTP
	url := e.endpoint + "/text/analytics/v3.1/entities/recognition/general"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Adăugăm header-ele necesare
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Ocp-Apim-Subscription-Key", e.apiKey)

	// Trimitem cererea
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Citim răspunsul
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	//fmt.Printf("API Response: %s\n", string(body)) // Log răspunsul brut

	// Verificăm codul de răspuns
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error: %s - %s", resp.Status, string(body))
	}

	// Deserializăm răspunsul
	var apiResp Response
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Verificăm dacă avem erori
	if len(apiResp.Errors) > 0 {
		return nil, fmt.Errorf("API returned error: %s", apiResp.Errors[0].Message)
	}

	// Verificăm dacă avem documente în răspuns
	if len(apiResp.Documents) == 0 {
		return nil, fmt.Errorf("API returned no documents")
	}

	// Construim rezultatul
	result := &models.ProcessResult{
		Success:  true,
		Entities: make([]models.Entity, 0),
	}

	// Adăugăm entitățile la rezultat
	for _, entity := range apiResp.Documents[0].Entities {
		matches := make([]string, 0)
		for _, match := range entity.Matches {
			matches = append(matches, match.Text)
		}

		result.Entities = append(result.Entities, models.Entity{
			Name:       entity.Name,
			Category:   entity.Category,
			Confidence: entity.Confidence,
			Offset:     entity.Offset,
			Length:     entity.Length,
			SubType:    entity.SubType,
			Matches:    matches,
		})
	}

	return result, nil
}
