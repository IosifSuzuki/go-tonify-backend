package converter

import (
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/domain/account/model"
)

func ConvertModel2CompanyResponse(model *model.Company) *dto.Company {
	return &dto.Company{
		ID:          model.ID,
		Name:        &model.Name,
		Description: &model.Description,
	}
}
