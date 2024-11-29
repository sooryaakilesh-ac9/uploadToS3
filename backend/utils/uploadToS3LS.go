package utils

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

func UploadToS3LSImages(path string, fileName string) error {
	// Initialize AWS session for LocalStack
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String("http://localhost:4566"),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
	})
	if err != nil {
		return fmt.Errorf("unable to create session: %v", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Define the S3 bucket name and folder path
	bucketName := os.Getenv("S3LS_BUCKET_NAME")
	folderPath := os.Getenv("S3LS_IMAGES_DIR_PATH")

	// Create S3 bucket (if it doesn't exist)
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		fmt.Println("Bucket may already exist, skipping creation")
	}

	// Read the image file (example: "image.jpg")
	imagePath := path // Path to your image file
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("unable to open file %v: %v", imagePath, err)
	}
	defer file.Close()

	// Read the file into a byte slice
	imageData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("unable to read file %v: %v", imagePath, err)
	}

	// Upload the image to the /images folder in S3
	objectKey := folderPath + fileName // Path in the bucket

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(imageData),
	})
	if err != nil {
		return fmt.Errorf("failed to upload image to S3: %v", err)
	}

	fmt.Printf("Successfully uploaded %s to s3://%s/%s\n", imagePath, bucketName, objectKey)
	return nil
}

func UploadQuotesMetadataToS3LS(path string, fileName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(os.Getenv("S3LS_REGION")),
		Endpoint:         aws.String(os.Getenv("S3LS_ENDPOINT")),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(os.Getenv("S3LS_ID"), os.Getenv("S3LS_SECRET"), os.Getenv("S3LS_TOKEN")),
	})
	if err != nil {
		return fmt.Errorf("unable to create session: %v", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Define the S3 bucket name and folder path
	bucketName := os.Getenv("S3LS_BUCKET_NAME")
	folderPath := os.Getenv("S3LS_QUOTES_DIR_PATH")

	// Create S3 bucket (if it doesn't exist)
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		fmt.Println("Bucket may already exist, skipping creation")
	}

	// Read the image file (example: "image.jpg")
	filePath := path // Path to your image file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file %v: %v", filePath, err)
	}
	defer file.Close()

	// Read the file into a byte slice
	imageData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("unable to read file %v: %v", filePath, err)
	}

	// Upload the image to the /images folder in S3
	objectKey := folderPath + fileName // Path in the bucket

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(imageData),
	})
	if err != nil {
		return fmt.Errorf("failed to upload image to S3: %v", err)
	}

	fmt.Printf("Successfully uploaded %s to s3://%s/%s\n", filePath, bucketName, objectKey)
	return nil
}

func UploadImagesMetadataToS3LS(path string, fileName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(os.Getenv("S3LS_REGION")),
		Endpoint:         aws.String(os.Getenv("S3LS_ENDPOINT")),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
	})
	if err != nil {
		return fmt.Errorf("unable to create session: %v", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Define the S3 bucket name and folder path
	bucketName := os.Getenv("S3LS_BUCKET_NAME")
	folderPath := os.Getenv("S3LS_IMAGES_DIR_PATH")

	// Create S3 bucket (if it doesn't exist)
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		fmt.Println("Bucket may already exist, skipping creation")
	}

	// Read the image file (example: "image.jpg")
	filePath := path // Path to your image file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file %v: %v", filePath, err)
	}
	defer file.Close()

	// Read the file into a byte slice
	imageData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("unable to read file %v: %v", filePath, err)
	}

	// Upload the image to the /images folder in S3
	objectKey := folderPath + fileName // Path in the bucket

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(imageData),
	})
	if err != nil {
		return fmt.Errorf("failed to upload image to S3: %v", err)
	}

	fmt.Printf("Successfully uploaded %s to s3://%s/%s\n", filePath, bucketName, objectKey)
	return nil
}