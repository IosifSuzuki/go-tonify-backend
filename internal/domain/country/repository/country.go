package repository

import (
	"encoding/json"
	"go-tonify-backend/internal/domain/entity"
	"io"
	"os"
)

type Country interface {
	GetAll() ([]entity.Country, error)
}

type country struct {
}

func NewCountry() Country {
	return &country{}
}

func (c *country) GetAll() ([]entity.Country, error) {
	file, err := os.Open("json/countries.json")
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var countries = make([]entity.Country, 0)
	if err := json.Unmarshal(data, &countries); err != nil {
		return nil, err
	}
	return countries, nil
}
