package converter

import (
	"go-tonify-backend/internal/domain/category/model"
	"go-tonify-backend/internal/domain/entity"
)

func ConvertEntity2CategoryModel(categoryEntity *entity.Category) *model.Category {
	return &model.Category{
		ID:    categoryEntity.ID,
		Title: categoryEntity.Title,
	}
}

func ConvertEntities2CategoriesModel(categoryEntity []entity.Category) []model.Category {
	var categories = make([]model.Category, 0, len(categoryEntity))
	for _, categoryEntity := range categoryEntity {
		categories = append(categories, *ConvertEntity2CategoryModel(&categoryEntity))
	}
	return categories
}
