package converter

import (
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/domain/country/model"
)

func ConvertModel2CountryResponse(country model.Country) *dto.Country {
	return &dto.Country{
		Name: country.Name,
		Code: country.Code,
	}
}
