package utils

import "path/filepath"

func ExtFromFileName(fileName string) string {
	return filepath.Ext(fileName)
}
