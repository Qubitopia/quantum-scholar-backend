package database

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Global S3 client for Cloudflare R2

var s3Client *s3.Client
var presignClient *s3.PresignClient

func InitR2Client() error {
	// Custom AWS config for R2
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(OBJ_ACCESS_KEY_ID, OBJ_SECRET_ACCESS_KEY, "")),
		config.WithRegion("auto"),
	)

	if err != nil {
		log.Fatal(err)
	}

	s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", OBJ_ACCOUNT_ID))
	})

	presignClient = s3.NewPresignClient(s3Client)

	return nil
}

func UploadObject(objectKey, contentType string, data []byte) error {
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &OBJ_BUCKET,
		Key:         &objectKey,
		Body:        bytes.NewReader(data),
		ContentType: &contentType,
	})

	return err
}

func GetPresignedURL(objectKey string, expiresIn time.Duration) (string, error) {
	req, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &OBJ_BUCKET,
		Key:    &objectKey,
	}, s3.WithPresignExpires(expiresIn))
	if err != nil {
		return "", err
	}

	return req.URL, nil
}
