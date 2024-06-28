package filesystem

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/JubaerHossain/rootx/pkg/core/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type FileUploadService struct {
	Config *config.Config
}

func NewFileUploadService(cfg *config.Config) *FileUploadService {
	return &FileUploadService{Config: cfg}
}

func FileUpload(r *http.Request, formKey string, cfg *config.Config, folder string) (string, error) {
	file, handler, err := r.FormFile(formKey)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Determine storage destination based on app config
	var filePath string
	switch cfg.StorageDisk {
	case "s3":
		filePath = folder + "/" + handler.Filename
		return uploadToS3(file, cfg, filePath)
	case "local":
		return saveToLocal(file, cfg, folder, handler.Filename)
	default:
		return "", errors.New("storage disk not supported")
	}
}

// uploadToS3 uploads file to AWS S3
func uploadToS3(file multipart.File, cfg *config.Config, filePath string) (string, error) {
	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AwsRegion),
		Credentials: credentials.NewStaticCredentials(cfg.AwsAccessKey, cfg.AwsSecretKey, ""),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %v", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %v", err)
	}

	// Upload file to S3 bucket
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(cfg.AwsBucket),
		Key:         aws.String(filePath),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(http.DetectContentType(fileBytes)),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// Construct public URL for the uploaded file
	fileURL := fmt.Sprintf("%s/%s/%s", cfg.AwsEndpoint, cfg.AwsBucket, filePath)
	return fileURL, nil
}

// saveToLocal saves file to local disk
func saveToLocal(file multipart.File, cfg *config.Config, folder string, filename string) (string, error) {
	// Define the root directory path
	rootDir := cfg.StoragePath

	// Create full file path
	filePath := filepath.Join(rootDir, folder, filename)

	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directories: %v", err)
	}

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %v", err)
	}

	// Write file to local disk
	if err := os.WriteFile(filePath, fileBytes, 0644); err != nil {
		return "", fmt.Errorf("failed to save file to local disk: %v", err)
	}

	// Return local file path with domain
	fileURL := fmt.Sprintf("%s/uploads/%s/%s", cfg.Domain, folder, filename)
	return fileURL, nil
}

// DeleteFromS3 deletes file from AWS S3
func DeleteFromS3(filePath string, cfg *config.Config) error {
	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AwsRegion),
		Credentials: credentials.NewStaticCredentials(cfg.AwsAccessKey, cfg.AwsSecretKey, ""),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Delete object from S3 bucket
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(cfg.AwsBucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %v", err)
	}

	return nil
}

// Example usage to delete an image (assuming you have the filePath stored)
func DeleteImage(filePath string, cfg *config.Config) error {
	var err error
	switch cfg.StorageDisk {
	case "s3":
		err = DeleteFromS3(filePath, cfg)
	case "local":
		err = DeleteFromLocal(filePath)
	default:
		err = errors.New("storage disk not supported")
	}
	return err
}

// DeleteFromLocal deletes file from local disk
func DeleteFromLocal(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file from local disk: %v", err)
	}
	return nil
}
