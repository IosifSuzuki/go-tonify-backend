package converter

import (
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/domain/account/model"
)

func ConvertModel2PairTokenResponse(pairTokenModel model.PairToken) *dto.PairToken {
	return &dto.PairToken{
		Access:  pairTokenModel.Access,
		Refresh: pairTokenModel.Refresh,
	}
}
