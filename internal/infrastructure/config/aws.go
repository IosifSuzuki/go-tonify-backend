package config

import (
	"go-tonify-backend/internal/domain/entity"
	"os"
	"sync"
)

type AWS struct {
	AttachmentBucket string
	Region           string
}

var (
	awsInstance *AWS
	awsErr      error
	awsOnce     sync.Once
)

func GetAWS() (*AWS, error) {
	awsOnce.Do(func() {
		var (
			instance AWS
			ok       bool
		)
		instance.AttachmentBucket, ok = os.LookupEnv("S3_ATTACHMENT_BUCKET")
		if !ok {
			awsErr = entity.NilError
			return
		}
		instance.Region, ok = os.LookupEnv("AWS_REGION")
		if !ok {
			awsErr = entity.NilError
			return
		}
		awsInstance = &instance
	})
	return awsInstance, awsErr
}
