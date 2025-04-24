package config

import (
	"errors"
	"os"
)

// Config structura pentru configurarea aplicației
type Config struct {
	// Azure Blob Storage
	AzureBlobAccountName   string
	AzureBlobAccountKey    string
	AzureBlobContainerName string

	// Azure SQL Database
	AzureSQLServerName string
	AzureSQLDBName     string
	AzureSQLUsername   string
	AzureSQLPassword   string

	// Azure Cognitive Services
	AzureCognitiveServicesEndpoint string
	AzureCognitiveServicesKey      string
}

// LoadConfig încarcă configurația din variabilele de mediu
func LoadConfig() (*Config, error) {
	config := &Config{
		// Azure Blob Storage
		AzureBlobAccountName:   os.Getenv("AZURE_BLOB_ACCOUNT_NAME"),
		AzureBlobAccountKey:    os.Getenv("AZURE_BLOB_ACCOUNT_KEY"),
		AzureBlobContainerName: os.Getenv("AZURE_BLOB_CONTAINER_NAME"),

		// Azure SQL Database
		AzureSQLServerName: os.Getenv("AZURE_SQL_SERVER_NAME"),
		AzureSQLDBName:     os.Getenv("AZURE_SQL_DB_NAME"),
		AzureSQLUsername:   os.Getenv("AZURE_SQL_USERNAME"),
		AzureSQLPassword:   os.Getenv("AZURE_SQL_PASSWORD"),

		// Azure Cognitive Services
		AzureCognitiveServicesEndpoint: os.Getenv("AZURE_COGNITIVE_SERVICES_ENDPOINT"),
		AzureCognitiveServicesKey:      os.Getenv("AZURE_COGNITIVE_SERVICES_KEY"),
	}

	// Verificare configuraţie minimă necesară
	if config.AzureBlobAccountName == "" || config.AzureBlobAccountKey == "" || config.AzureBlobContainerName == "" {
		return nil, errors.New("missing Azure Blob Storage configuration")
	}

	if config.AzureSQLServerName == "" || config.AzureSQLDBName == "" || config.AzureSQLUsername == "" || config.AzureSQLPassword == "" {
		return nil, errors.New("missing Azure SQL Database configuration")
	}

	if config.AzureCognitiveServicesEndpoint == "" || config.AzureCognitiveServicesKey == "" {
		return nil, errors.New("missing Azure Cognitive Services configuration")
	}

	return config, nil
}
