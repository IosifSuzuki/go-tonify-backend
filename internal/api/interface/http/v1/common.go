package v1

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/api/interface/http/v1/converter"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain/country/usecase"
	"go-tonify-backend/pkg/logger"
	"net/http"
)

type CommonHandler struct {
	container      container.Container
	countryUsecase usecase.Country
}

func NewCommonHandler(container container.Container, countryUsecase usecase.Country) *CommonHandler {
	return &CommonHandler{
		container:      container,
		countryUsecase: countryUsecase,
	}
}

// Ping godoc
//
//	@Summary		Ping to server
//	@Description	Simple request to check if the server is alive
//	@Tags			common
//	@Produce		json
//	@Success		200	{object}	dto.Response{response=string}	"returns 'ok' if the server is alive"
//	@Router			/v1/common/ping [get]
func (c *CommonHandler) Ping(ctx *gin.Context) {
	successResponse(ctx, http.StatusOK, "ok")
}

// Countries godoc
//
//	@Summary		Get countries
//	@Description	Retries all available countries
//	@Tags			common
//	@Produce		json
//	@Success		200	{object}	dto.Response{response=[]dto.Country}	"countries"
//	@Success		500	{object}	dto.Response{response=dto.Empty}		"detailed error message"
//	@Router			/v1/common/countries [get]
func (c *CommonHandler) Countries(ctx *gin.Context) {
	log := c.container.GetLogger()
	countryModels, err := c.countryUsecase.GetCountries()
	if err != nil {
		log.Error("fail to get all countries", logger.FError(err))
		failResponse(ctx, http.StatusInternalServerError, dto.InternalServerError, err)
		return
	}
	countries := make([]dto.Country, 0, len(countryModels))
	for _, countryModel := range countryModels {
		country := converter.ConvertModel2CountryResponse(countryModel)
		countries = append(countries, *country)
	}
	successResponse(ctx, http.StatusOK, countries)
}
