package utils

import (
	"go-tonify-backend/internal/domain/entity"
	"path/filepath"
)

func ExtFromFileName(fileName string) (*string, error) {
	ext := filepath.Ext(fileName)
	if len(ext) == 0 {
		return nil, entity.EmptyValueError
	}
	return &ext, nil
}
