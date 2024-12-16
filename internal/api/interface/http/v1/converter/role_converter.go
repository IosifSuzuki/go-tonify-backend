package converter

import (
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/domain/account/model"
)

func ConvertDto2RoleModel(role dto.Role) model.Role {
	switch role {
	case dto.ClientRole:
		return model.ClientRole
	case dto.FreelancerRole:
		return model.FreelancerRole
	default:
		return model.UnknownRole
	}
}
