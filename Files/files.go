package Files

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jonathanpatta/apartmentservices/Settings"
	"strings"
)

type Images struct {
	Data     string `json:"data,omitempty"`
	FileName string `json:"file_name,omitempty"`
	FileType string `json:"file_type,omitempty"`
	Bytes    []byte `json:"bytes,omitempty"`
}

type S3FileService struct {
	cli        *s3.Client
	s3Settings *Settings.S3Settings
	region     string
}

func NewS3FileService(settings *Settings.Settings) (*S3FileService, error) {
	return &S3FileService{
		cli:        settings.S3Settings.Cli,
		s3Settings: settings.S3Settings,
		region:     settings.Region,
	}, nil
}

func (s *S3FileService) UploadImages(imgs []Images) ([]string, error) {
	var urls []string
	for _, img := range imgs {
		key := "static/images/" + img.FileName
		imgBytes, err := ProcessImage(img.Bytes)
		if err != nil {
			return nil, err
		}
		_, err = s.cli.PutObject(context.Background(), &s3.PutObjectInput{
			Key:         aws.String(key),
			Bucket:      aws.String(s.s3Settings.BucketName),
			Body:        bytes.NewReader(imgBytes),
			ContentType: aws.String(img.FileType),
		})
		if err != nil {
			return urls, err
		}
		key = strings.Replace(key, " ", "+", -1)
		newUrl := fmt.Sprintf("https://%v.s3-%v.amazonaws.com/%v", s.s3Settings.BucketName, s.region, key)
		urls = append(urls, newUrl)
	}
	return urls, nil
}

func (s *S3FileService) DeleteImages() ([]string, error) {
	return nil, nil
}
