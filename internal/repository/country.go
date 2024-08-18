package repository

import (
	"encoding/json"
	"go-tonify-backend/internal/domain"
	"io"
	"os"
)

type CountryRepository interface {
	FetchAllCountries() ([]domain.Country, error)
}

type countryRepository struct {
}

func NewCountryRepository() CountryRepository {
	return &countryRepository{}
}

func (c *countryRepository) FetchAllCountries() ([]domain.Country, error) {
	file, err := os.Open("json/countries.json")
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var countries = make([]domain.Country, 0)
	if err := json.Unmarshal(data, &countries); err != nil {
		return nil, err
	}
	return countries, nil
}
