package converter

import (
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/domain/account/model"
)

func ConvertModel2TagResponse(tagModel *model.Tag) *dto.Tag {
	return &dto.Tag{
		ID:    tagModel.ID,
		Title: tagModel.Title,
	}
}

func ConvertModels2TagsResponse(tagModels []model.Tag) []dto.Tag {
	var tags = make([]dto.Tag, 0, len(tagModels))
	for _, tagModel := range tagModels {
		tags = append(tags, *ConvertModel2TagResponse(&tagModel))
	}
	return tags
}
