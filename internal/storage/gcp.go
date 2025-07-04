package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Shubham-Thakur06/go-streaming-platform/internal/config"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCPProvider struct {
	client     *storage.Client
	bucketName string
}

func NewGCPProvider(cfg config.GCPConfig) (*GCPProvider, error) {
	ctx := context.Background()

	var client *storage.Client
	var err error

	if cfg.CredentialsFile != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(cfg.CredentialsFile))
	} else {
		client, err = storage.NewClient(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create GCP storage client: %w", err)
	}

	return &GCPProvider{
		client:     client,
		bucketName: cfg.BucketName,
	}, nil
}

func (p *GCPProvider) UploadFile(file io.Reader, filename string, contentType string) (string, error) {
	ctx := context.Background()
	bucket := p.client.Bucket(p.bucketName)
	obj := bucket.Object(filename)

	writer := obj.NewWriter(ctx)
	writer.ContentType = contentType

	if _, err := io.Copy(writer, file); err != nil {
		return "", fmt.Errorf("failed to copy file to GCP: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	return p.GetFileURL(filename)
}

func (p *GCPProvider) DeleteFile(filename string) error {
	ctx := context.Background()
	bucket := p.client.Bucket(p.bucketName)
	obj := bucket.Object(filename)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete file from GCP: %w", err)
	}
	return nil
}

func (p *GCPProvider) GetFileURL(filename string) (string, error) {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", p.bucketName, filename), nil
}

func (p *GCPProvider) GeneratePresignedURL(filename string, expiresIn int64) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(time.Duration(expiresIn) * time.Second),
	}

	url, err := storage.SignedURL(p.bucketName, filename, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}
