package s3

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Service struct {
	session *session.Session
}

func NewS3Service() (*S3Service, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(os.Getenv("S3_REGION")),
		Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(os.Getenv("S3_ID"), os.Getenv("S3_SECRET"), os.Getenv("S3_TOKEN")),
	})
	if err != nil {
		return nil, err
	}
	return &S3Service{session: sess}, nil
}

func (s *S3Service) UploadImage(filePath, fileName string) error {
	return s.uploadFile(filePath, fileName, os.Getenv("S3_IMAGES_DIR_PATH"))
}

func (s *S3Service) UploadMetadata(filePath, fileName string) error {
	return s.uploadFile(filePath, fileName, "")
}

func (s *S3Service) uploadFile(filePath, fileName, dirPath string) error {
	svc := s3.New(s.session)
	bucketName := os.Getenv("S3_BUCKET_NAME")

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file %v: %v", filePath, err)
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("unable to read file %v: %v", filePath, err)
	}

	objectKey := fileName
	if dirPath != "" {
		objectKey = dirPath + fileName
	}

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(fileData),
	})
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %v", err)
	}

	return nil
} 