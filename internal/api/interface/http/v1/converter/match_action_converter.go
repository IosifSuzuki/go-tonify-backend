package converter

import (
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/domain/account/model"
)

func ConvertDto2MatchActionModel(matchAction dto.MatchAction) model.MatchAction {
	switch matchAction {
	case dto.LikeMatchAction:
		return model.LikeMatchAction
	case dto.DislikeMatchAction:
		return model.DislikeMatchAction
	default:
		return model.UnknownMatchAction
	}
}
