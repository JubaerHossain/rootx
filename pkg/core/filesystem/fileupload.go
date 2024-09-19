package filesystem

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JubaerHossain/rootx/pkg/core/config"

	rootUtils "github.com/JubaerHossain/rootx/pkg/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type FileMetadata struct {
	Name         string
	OriginalName string
	RootPath     string
	Extension    string
	Path         string
	Type         string
	URL          string
}

type FileUploadService struct {
	Config *config.Config
}

func NewFileUploadService(cfg *config.Config) *FileUploadService {
	return &FileUploadService{Config: cfg}
}

func (s *FileUploadService) FileUpload(r *http.Request, formKey string, folder string) (*FileMetadata, error) {
	file, handler, err := r.FormFile(formKey)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Generate a unique file name
	uniqueFileName := s.generateUniqueFileName(handler.Filename)

	// Determine storage destination based on app config
	if strings.TrimSpace(s.Config.StorageDisk) == "s3" {
		filePath := folder + "/" + uniqueFileName
		return s.uploadToS3(file, s.Config, filePath, handler.Filename)
	} else {
		return s.saveToLocal(file, s.Config, folder, uniqueFileName)
	}
}

func (s *FileUploadService) generateUniqueFileName(originalName string) string {
	// Generate a timestamp

	// Extract the file extension
	ext := filepath.Ext(originalName)
	name := originalName[:len(originalName)-len(ext)]

	// Create a new file name by appending the unique before the extension
	uniqueNumber, err := rootUtils.GenerateUniqueNumber(8)
	if err != nil {
		uniqueNumber = time.Now().Format("20060102150405")
	}
	newFileName := fmt.Sprintf("%s_%s%s", name, uniqueNumber, ext)

	return newFileName
}

// uploadToS3 uploads file to AWS S3 and returns file metadata
func (s *FileUploadService) uploadToS3(file multipart.File, cfg *config.Config, filePath string, originalName string) (*FileMetadata, error) {
	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AwsRegion),
		Credentials: credentials.NewStaticCredentials(cfg.AwsAccessKey, cfg.AwsSecretKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %v", err)
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
		return nil, fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// Construct public URL for the uploaded file
	fileURL := fmt.Sprintf("%s/%s/%s", cfg.AwsEndpoint, cfg.AwsBucket, filePath)

	// Return file metadata
	metadata := &FileMetadata{
		Name:         filepath.Base(filePath),
		OriginalName: originalName,
		RootPath:     cfg.AwsBucket,
		Extension:    filepath.Ext(filePath),
		Path:         filePath,
		Type:         http.DetectContentType(fileBytes),
		URL:          fileURL,
	}
	return metadata, nil
}

// saveToLocal saves file to local disk and returns file metadata
func (s *FileUploadService) saveToLocal(file multipart.File, cfg *config.Config, folder string, filename string) (*FileMetadata, error) {
	// Define the root directory path
	rootDir := cfg.StoragePath

	// Create full file path
	filePath := filepath.Join(rootDir, folder, filename)

	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create directories: %v", err)
	}

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %v", err)
	}

	// Write file to local disk
	if err := os.WriteFile(filePath, fileBytes, 0644); err != nil {
		return nil, fmt.Errorf("failed to save file to local disk: %v", err)
	}

	// Return file metadata
	fileURL := fmt.Sprintf("%s/uploads/%s/%s", cfg.Domain, folder, filename)
	metadata := &FileMetadata{
		Name:         filename,
		OriginalName: filename,
		RootPath:     rootDir,
		Extension:    filepath.Ext(filename),
		Path:         filePath,
		Type:         http.DetectContentType(fileBytes),
		URL:          fileURL,
	}
	return metadata, nil
}

// DeleteFromS3 deletes file from AWS S3
func (s *FileUploadService) DeleteFromS3(filePath string, cfg *config.Config) error {
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
func (s *FileUploadService) DeleteImage(filePath string) error {
	var err error
	if strings.TrimSpace(s.Config.StorageDisk) == "s3" {
		err = s.DeleteFromS3(filePath, s.Config)
	} else {
		err = s.DeleteFromLocal(filePath)
	}
	return err
}

// DeleteFromLocal deletes file from local disk
func (s *FileUploadService) DeleteFromLocal(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file from local disk: %v", err)
	}
	return nil
}
