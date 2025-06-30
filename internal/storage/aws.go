package storage

import (
	"fmt"
	"io"

	"github.com/Shubham-Thakur06/go-streaming-platform/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSProvider struct {
	s3Client   *s3.S3
	uploader   *s3manager.Uploader
	bucketName string
}

func NewAWSProvider(cfg config.AWSConfig) (*AWSProvider, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	s3Client := s3.New(sess)
	uploader := s3manager.NewUploader(sess)

	return &AWSProvider{
		s3Client:   s3Client,
		uploader:   uploader,
		bucketName: cfg.BucketName,
	}, nil
}

func (p *AWSProvider) UploadFile(file io.Reader, filename string, contentType string) (string, error) {
	return "", nil
}
func (p *AWSProvider) DeleteFile(filename string) error {
	return nil
}
func (p *AWSProvider) GetFileURL(filename string) (string, error) {
	return "", nil
}
func (p *AWSProvider) GeneratePresignedURL(filename string, expiresIn int64) (string, error) {
	return "", nil
}
