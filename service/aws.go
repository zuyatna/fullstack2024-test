package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	// "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	// "github.com/aws/aws-sdk-go-v2/service/s3" // error go: github.com/aws/aws-sdk-go-v2/service/internal/checksum imports hash/crc64: package hash/crc64 is not in std (C:\Program Files\Go\src\hash\crc64)
)

type S3Service struct {
	// s3Client *s3.Client
	bucket string
}

func NewS3Service() (*S3Service, error) {
	awsRegion := os.Getenv("AWS_REGION")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")
	bucketName := os.Getenv("AWS_S3_BUCKET")

	if awsRegion == "" || awsAccessKey == "" || awsSecretKey == "" || bucketName == "" {
		return nil, fmt.Errorf("AWS configuration missing")
	}

	_, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			awsAccessKey, awsSecretKey, "",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return nil, nil
}

func (s *S3Service) UploadFile(file *multipart.FileHeader, clientID int) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	buffer := make([]byte, file.Size)
	if _, err = src.Read(buffer); err != nil {
		return "", err
	}

	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("clients/%d/logo%s", clientID, fileExt)

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// NOTE: error import aws package hash/crc64
	// link: https://go.dev/src/hash/crc64/crc64.go

	region := os.Getenv("AWS_REGION")
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, region, fileName)

	return fileURL, nil
}
