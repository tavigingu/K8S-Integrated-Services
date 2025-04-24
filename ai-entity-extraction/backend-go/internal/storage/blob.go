package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"backend-go/internal/config"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// BlobStorage reprezintă un serviciu pentru interacțiunea cu Azure Blob Storage
type BlobStorage struct {
	containerURL azblob.ContainerURL
	cfg          *config.Config
}

// NewBlobStorage creează o nouă instanță a serviciului BlobStorage
func NewBlobStorage(cfg *config.Config) (*BlobStorage, error) {
	// Creează o credențial folosind cheia contului de storage
	credential, err := azblob.NewSharedKeyCredential(cfg.AzureBlobAccountName, cfg.AzureBlobAccountKey)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %v", err)
	}

	// Creează o pipeline folosind credențialul
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Construiește URL-ul pentru container-ul de blob-uri
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", cfg.AzureBlobAccountName, cfg.AzureBlobContainerName))
	containerURL := azblob.NewContainerURL(*URL, pipeline)

	// Creează container-ul dacă nu există
	ctx := context.Background()
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		// Verificăm dacă eroarea este că container-ul există deja, ceea ce este OK
		if stgErr, ok := err.(azblob.StorageError); ok {
			if stgErr.ServiceCode() == azblob.ServiceCodeContainerAlreadyExists {
				err = nil
			}
		}
		if err != nil {
			return nil, fmt.Errorf("failed to create container: %v", err)
		}
	}

	return &BlobStorage{
		containerURL: containerURL,
		cfg:          cfg,
	}, nil
}

// UploadFile încarcă un fișier în Azure Blob Storage
func (bs *BlobStorage) UploadFile(ctx context.Context, filename string, data io.Reader, contentType string) (string, error) {
	// Creează un blob URL pentru fișier
	blobURL := bs.containerURL.NewBlockBlobURL(filename)

	// Încarcă fișierul
	// Notă: În versiunea mai nouă, tipul s-a schimbat la azblob.UploadStreamToBlockBlobOptions
	// în loc de azblob.UploadToBlockBlobOptions
	uploadOptions := azblob.UploadStreamToBlockBlobOptions{
		BufferSize: 4 * 1024 * 1024, // Dimensiune bloc: 4MB
		MaxBuffers: 16,              // Încărcare paralelă
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: contentType,
		},
		Metadata: azblob.Metadata{
			"uploaded": time.Now().Format(time.RFC3339),
		},
	}

	_, err := azblob.UploadStreamToBlockBlob(ctx, data, blobURL, uploadOptions)
	if err != nil {
		return "", fmt.Errorf("failed to upload blob: %v", err)
	}

	urlPtr := blobURL.URL()
	blobURLStr := urlPtr.String()

	return blobURLStr, nil
}

// GetFile descarcă un fișier din Azure Blob Storage
func (bs *BlobStorage) GetFile(ctx context.Context, filename string) (io.ReadCloser, error) {
	// Creează un blob URL pentru fișier
	blobURL := bs.containerURL.NewBlockBlobURL(filename)

	// Descarcă blob-ul
	response, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download blob: %v", err)
	}

	// Returnează un reader pentru conținutul blob-ului
	return response.Body(azblob.RetryReaderOptions{}), nil
}

// DeleteFile șterge un fișier din Azure Blob Storage
func (bs *BlobStorage) DeleteFile(ctx context.Context, filename string) error {
	// Creează un blob URL pentru fișier
	blobURL := bs.containerURL.NewBlockBlobURL(filename)

	// Șterge blob-ul
	_, err := blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		return fmt.Errorf("failed to delete blob: %v", err)
	}

	return nil
}
