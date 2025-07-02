package storage

import (
	"fmt"
	"io"
	"time"

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
	_, err := p.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(p.bucketName),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return p.GetFileURL(filename)
}

func (p *AWSProvider) DeleteFile(filename string) error {
	_, err := p.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}
	return nil
}

func (p *AWSProvider) GetFileURL(filename string) (string, error) {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", p.bucketName, *p.s3Client.Config.Region, filename), nil
}

func (p *AWSProvider) GeneratePresignedURL(filename string, expiresIn int64) (string, error) {
	req, _ := p.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(filename),
	})

	url, err := req.Presign(time.Duration(expiresIn) * time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url, nil
}
