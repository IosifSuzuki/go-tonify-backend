package converter

import (
	"go-tonify-backend/internal/domain/account/model"
	"go-tonify-backend/internal/domain/entity"
)

func ConvertEntity2TagModel(tagEntity *entity.Tag) *model.Tag {
	return &model.Tag{
		ID:    tagEntity.ID,
		Title: tagEntity.Title,
	}
}

func ConvertEntities2TagModels(tagEntities []entity.Tag) []model.Tag {
	var tagModels = make([]model.Tag, 0, len(tagEntities))
	for _, tagEntity := range tagEntities {
		tagModels = append(tagModels, *ConvertEntity2TagModel(&tagEntity))
	}
	return tagModels
}
