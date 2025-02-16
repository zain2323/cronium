package files

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	AwsS3Bucket = "cronium-file-service"
	AwsS3Region = "us-east-2"
)

// S3Storage S3 is an implementation of the Storage interface which works with the AWS S3
type S3Storage struct {
	maxFileSize int64
	bucketName  string
	client      *s3.Client
	uploader    *manager.Uploader
	downloader  *manager.Downloader
}

// createClient initializes the S3 client by using the AWS ACCESS_KEY and SECRET_KEY
// configured through the AWS cli
// TODO: replace it with env vars
func createClient() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AwsS3Region))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	s3Client := s3.NewFromConfig(cfg)
	return s3Client, nil
}

// NewS3Storage creates a new S3 file storage with the given bucket name
// bucketName is the bucket that will be created with the provided name.
// maxSize is the max number of bytes that an S3 object can be
func NewS3Storage(maxSize int64) (*S3Storage, error) {
	client, err := createClient()
	if err != nil {
		return nil, err
	}
	return &S3Storage{
			bucketName:  AwsS3Bucket,
			maxFileSize: maxSize,
			client:      client,
			uploader:    manager.NewUploader(client),
			downloader:  manager.NewDownloader(client)},
		nil
}

// Save the contents of the writer as the object with the provided name
func (s *S3Storage) Save(objectName string, contents io.Reader) error {

	// Upload the file to S3
	_, err := s.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(AwsS3Bucket),
		Key:    aws.String(objectName),
		Body:   contents,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("File uploaded successfully")
	return nil
}

// Get reads the object from the bucket
func (s *S3Storage) Get(objectName string) (*os.File, error) {
	fp := s.fullPath(objectName)
	// get the directory and make sure it exists
	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// create the new file
	newFile, err := os.Create(fp)
	if err != nil {
		log.Println("sss")
		log.Println(err)
	}
	defer newFile.Close()
	log.Println("File retrieved successfully")

	_, err = s.downloader.Download(context.TODO(), newFile, &s3.GetObjectInput{
		Bucket: aws.String(AwsS3Bucket),
		Key:    aws.String(objectName),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return newFile, nil
}

// returns the absolute path
func (s *S3Storage) fullPath(path string) string {
	return filepath.Join(s.bucketName, path)
}
