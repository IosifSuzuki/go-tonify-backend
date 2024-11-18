package model

import (
	"go-tonify-backend/internal/utils"
	"io"
	"path/filepath"
)

type UploadFile struct {
	Name string
	Body io.Reader
}

func (u *UploadFile) Ext() *string {
	return utils.NewString(filepath.Ext(u.Name)[1:])
}
