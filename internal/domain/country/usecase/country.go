package usecase

import (
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain/country/model"
	"go-tonify-backend/internal/domain/country/repository"
	"go-tonify-backend/pkg/logger"
)

type Country interface {
	GetCountries() ([]model.Country, error)
}

type country struct {
	container         container.Container
	countryRepository repository.Country
}

func NewCountry(
	container container.Container,
	countryRepository repository.Country,
) Country {
	return &country{
		container:         container,
		countryRepository: countryRepository,
	}
}

func (c *country) GetCountries() ([]model.Country, error) {
	log := c.container.GetLogger()
	countries, err := c.countryRepository.GetAll()
	if err != nil {
		log.Error("fail to get all countries", logger.FError(err))
		return nil, err
	}
	countryModels := make([]model.Country, 0, len(countries))
	for _, country := range countries {
		countryModel := model.Country{
			Name: country.Name,
			Code: country.Code,
		}
		countryModels = append(countryModels, countryModel)
	}
	return countryModels, nil
}
