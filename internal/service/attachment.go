package service

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/model"
	"go-tonify-backend/internal/utils"
	"go-tonify-backend/pkg/logger"
)

type AttachmentService interface {
	UploadFile(file model.UploadFile) (*string, error)
}

type attachmentService struct {
	container container.Container
	session   *session.Session
}

func NewAttachmentService(container container.Container) AttachmentService {
	return &attachmentService{
		container: container,
		session:   session.Must(session.NewSession()),
	}
}

func (a *attachmentService) UploadFile(file model.UploadFile) (*string, error) {
	log := a.container.GetLogger()
	uploader := s3manager.NewUploader(a.session)
	name := uuid.NewString()
	ext := file.Ext()
	key := fmt.Sprintf("%s.%s", name, *ext)
	output, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: utils.NewString(a.container.GetS3Config().AttachmentBucket),
		Key:    utils.NewString(key),
		Body:   file.Body,
	})
	if err != nil {
		log.Error("fail to upload file on s3", logger.FError(err))
		return nil, err
	}
	log.Debug("get result from s3 bucket", logger.F("output", output))
	return &output.Location, nil
}
