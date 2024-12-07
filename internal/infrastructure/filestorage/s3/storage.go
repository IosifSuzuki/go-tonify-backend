package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/utils"
	"go-tonify-backend/pkg/logger"
	"mime/multipart"
)

type FileStorage struct {
	container  container.Container
	bucketName string
	sess       *session.Session
	s3Client   *s3.S3
}

func NewS3FileStorage(container container.Container) *FileStorage {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(container.GetAWSConfig().Region),
	}))
	return &FileStorage{
		container:  container,
		bucketName: container.GetAWSConfig().AttachmentBucket,
		sess:       sess,
		s3Client:   s3.New(sess),
	}
}

func (f *FileStorage) UploadFile(fileName string, file multipart.File) (*string, error) {
	log := f.container.GetLogger()
	uploader := s3manager.NewUploader(f.sess)
	output, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: utils.NewString(f.container.GetAWSConfig().AttachmentBucket),
		Key:    utils.NewString(fileName),
		Body:   file,
	})
	if err != nil {
		log.Error("fail to upload file on s3", logger.FError(err))
		return nil, err
	}
	log.Debug("get result from s3 bucket", logger.F("output", output))
	return &output.Location, nil
}

func (f *FileStorage) GetFileURL(fileName string) (*string, error) {
	return utils.NewString(
		fmt.Sprintf("https://%s.s3.amazonaws.com/%s", f.bucketName, fileName),
	), nil
}
func (f *FileStorage) DeleteFile(fileName string) error {
	log := f.container.GetLogger()
	log.Error("will remove object", logger.F("key", fileName))
	_, err := f.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: utils.NewString(f.bucketName),
		Key:    utils.NewString(fileName),
	})
	if err != nil {
		log.Error("fail to delete object", logger.F("key", fileName))
		return err
	}
	return nil
}
