package v1

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/api/interface/http/v1/converter"
	"go-tonify-backend/internal/api/interface/http/validator"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain/category/usecase"
	"go-tonify-backend/pkg/logger"
	"net/http"
)

type CategoryHandler struct {
	container       container.Container
	validation      validator.HttpValidator
	categoryUsecase usecase.Category
}

func NewCategoryHandler(
	container container.Container,
	validation validator.HttpValidator,
	categoryUsecase usecase.Category,
) *CategoryHandler {
	return &CategoryHandler{
		container:       container,
		validation:      validation,
		categoryUsecase: categoryUsecase,
	}
}

// GetAll godoc
//
//	@Summary		Get categories
//	@Description	Retries all categories
//	@Tags			category
//	@Produce		json
//	@Success		200	{object}	dto.Response{response=dto.Pagination{data=[]dto.Category}}	"categories"
//	@Failure		400	{object}	dto.Response{response=dto.Empty}		"detailed error message"
//	@Failure		401	{object}	dto.Response{response=dto.Empty}		"the authorization token is invalid/expired/missing"
//	@Failure		500	{object}	dto.Response{response=dto.Empty}		"detailed error message"
//	@Router			/v1/category/all [get]
//	@Security		ApiKeyAuth
func (c *CategoryHandler) GetAll(ctx *gin.Context) {
	log := c.container.GetLogger()
	var getCategories dto.GetCategories
	if err := ctx.ShouldBindQuery(&getCategories); err != nil {
		log.Error("fail to bind get list task", logger.FError(err))
		badRequestResponse(ctx, c.validation, dto.BadRequestError, err)
		return
	}
	paginationModel, err := c.categoryUsecase.GetAllCategories(ctx, getCategories.Offset, getCategories.Limit)
	if err != nil {
		log.Error(
			"fail to get categories with parameters",
			logger.F("offset", getCategories.Offset),
			logger.F("limit", getCategories.Limit),
		)
		failResponse(ctx, http.StatusInternalServerError, err, dto.FailProcessRequestError)
		return
	}
	categories := converter.ConvertModels2CategoriesResponse(paginationModel.Data)
	pagination := dto.Pagination{
		Offset: paginationModel.Offset,
		Limit:  paginationModel.Limit,
		Total:  paginationModel.Total,
		Data:   categories,
	}
	successResponse(ctx, http.StatusOK, pagination)
}
