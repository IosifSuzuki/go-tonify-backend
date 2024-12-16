package converter

import (
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/domain/account/model"
)

func ConvertModel2MatchActionResponse(matchAction model.MatchResult) dto.MatchResult {
	switch matchAction {
	case model.LikeMatchResult:
		return dto.LikeMatchResult
	case model.DislikeMatchResult:
		return dto.DislikeMatchResult
	case model.MatchAccountMatchResult:
		return dto.MatchAccountMatchResult
	default:
		return dto.ErrorMatchResult
	}
}
