package filestorage

import "mime/multipart"

type FileStorage interface {
	UploadFile(string, multipart.File) (*string, error)
	GetFileURL(string) (*string, error)
	DeleteFile(string) error
}
