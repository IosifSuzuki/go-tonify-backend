package converter

import (
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/domain/account/model"
)

func ConvertModel2AttachmentResponse(model *model.Attachment) *dto.Attachment {
	return &dto.Attachment{
		ID:   &model.ID,
		Name: &model.Name,
		Path: &model.Path,
	}
}
