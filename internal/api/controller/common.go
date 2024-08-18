package controller

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/model"
	"go-tonify-backend/internal/repository"
	"net/http"
)

type CommonController struct {
	Container         container.Container
	CountryRepository repository.CountryRepository
}

// Countries godoc
//
//	@Summary		list of countries
//	@Description	get a list of countries paired with name and code
//	@Tags			common
//	@Produce		json
//	@Success		200	{array}	model.Country	"countries"
//	@Failure		500	"internal error"
//	@Router			/common/countries [get]
func (c *CommonController) Countries(ctx *gin.Context) {
	allCountries, err := c.CountryRepository.FetchAllCountries()
	if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
	}
	countries := make([]model.Country, 0, len(allCountries))
	for _, country := range allCountries {
		countries = append(countries, model.Country{
			Name: country.Name,
			Code: country.Code,
		})
	}
	sendResponseWithStatus(ctx, countries, http.StatusOK)
}
