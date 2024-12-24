package converter

import (
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/domain/category/model"
)

func ConvertModel2CategoryResponse(model *model.Category) *dto.Category {
	return &dto.Category{
		ID:    model.ID,
		Title: model.Title,
	}
}

func ConvertModels2CategoriesResponse(models []model.Category) []dto.Category {
	var categories = make([]dto.Category, 0, len(models))
	for _, model := range models {
		categories = append(categories, *ConvertModel2CategoryResponse(&model))
	}
	return categories
}
