package v1

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/api/interface/http/v1/converter"
	"go-tonify-backend/internal/api/interface/http/validator"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain/account/model"
	"go-tonify-backend/internal/domain/account/usecase"
	"go-tonify-backend/pkg/logger"
	"net/http"
)

type MatchHandler struct {
	container    container.Container
	validation   validator.HttpValidator
	matchUsecase usecase.Match
}

func NewMatchHandler(
	container container.Container,
	validation validator.HttpValidator,
	matchUsecase usecase.Match,
) *MatchHandler {
	return &MatchHandler{
		container:    container,
		validation:   validation,
		matchUsecase: matchUsecase,
	}
}

// MatchableAccounts godoc
//
//	@Summary		Matchable accounts
//	@Description	Get matchable accounts: accounts that have not been liked, disliked, or were disliked a long time ago.
//	@Description	**Attention**: The rules may change from time to time. If you need more information about the endpoint, please contact API support
//	@Tags			match
//	@Param			Authorization	header		string					true	"account's access token"
//	@Param			request			body		dto.GetMatchAccounts	true	"matching accounts parameters"
//	@Produce		json
//	@Success		200	{object}	dto.Response{response=[]dto.Account}	"list of matchable accounts"
//	@Failure		400	{object}	dto.Response{response=dto.Empty}		"detailed error message"
//	@Failure		401	{object}	dto.Response{response=dto.Empty}		"the authorization token is invalid/expired/missing"
//	@Failure		410	{object}	dto.Response{response=dto.Empty}		"account does not exist or has been deleted"
//	@Failure		500	{object}	dto.Response{response=dto.Empty}		"detailed error message"
//	@Router			/v1/match/matchable/accounts [get]
//	@Security		ApiKeyAuth
func (m *MatchHandler) MatchableAccounts(ctx *gin.Context) {
	log := m.container.GetLogger()
	accountID, err := getAccountID(ctx)
	if err != nil {
		log.Error("fail to get account id", logger.FError(err))
		failResponse(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	var getMatchAccounts dto.GetMatchAccounts
	if err := ctx.ShouldBindJSON(&getMatchAccounts); err != nil {
		log.Error("fail to bind get match accounts", logger.FError(err))
		badRequestResponse(ctx, m.validation, dto.BadRequestError, err)
		return
	}
	accountModels, err := m.matchUsecase.MatchableAccounts(ctx, *accountID, getMatchAccounts.Limit)
	if err != nil {
		log.Error("fail to get matchable accounts", logger.FError(err))
		switch err {
		case model.EntityNotFoundError:
			failResponse(ctx, http.StatusGone, dto.ModelNotFoundError, err)
		default:
			failResponse(ctx, http.StatusInternalServerError, dto.FailProcessRequestError, err)
		}
		return
	}
	accounts := make([]dto.Account, 0, len(accountModels))
	for _, accountModel := range accountModels {
		account := converter.ConvertModel2AccountResponse(&accountModel)
		accounts = append(accounts, *account)
	}
	successResponse(ctx, http.StatusOK, accounts)
}

// MatchAction godoc
//
// @Summary		match action
// @Description  When a client performs a **like** action, the server can return one of the following responses:
// @Description  - **like**: The action was successful.
// @Description  - **error**: An error occurred while processing the request.
// @Description  - **match**: A mutual "like" was identified.
// @Description
// @Description  When a client performs a **dislike** action, the server can return:
// @Description  - **dislike**: The action was successful.
// @Description  - **error**: An error occurred while processing the request.
// @Tags			match
// @Param			Authorization	header		string					true	"account's access token"
// @Param			request			body		dto.PostMatchAction		true	"action match parameters"
// @Produce			json
// @Param        	action    		query     dto.MatchAction  true  "match action"
// @Success		200	{object}	dto.Response{response=dto.MatchResult}	"list of matching accounts"
// @Failure		400	{object}	dto.Response{response=dto.Empty}		"detailed error message"
// @Failure		401	{object}	dto.Response{response=dto.Empty}		"the authorization token is invalid/expired/missing"
// @Failure		500	{object}	dto.Response{response=dto.Empty}		"detailed error message"
// @Router			/v1/match/action/{action} [post]
// @Security		ApiKeyAuth
func (m *MatchHandler) MatchAction(ctx *gin.Context) {
	log := m.container.GetLogger()
	accountID, err := getAccountID(ctx)
	if err != nil {
		log.Error("fail to get account id", logger.FError(err))
		failResponse(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	var uriMatchAction dto.URIMatchAction
	if err := ctx.ShouldBindUri(&uriMatchAction); err != nil {
		log.Error("fail to bind uri match action", logger.FError(err))
		badRequestResponse(ctx, m.validation, dto.BadRequestError, err)
		return
	}
	var postMatchAction dto.PostMatchAction
	if err := ctx.ShouldBindJSON(&postMatchAction); err != nil {
		log.Error("fail to bind post match action", logger.FError(err))
		badRequestResponse(ctx, m.validation, dto.BadRequestError, err)
		return
	}
	matchActionModel := converter.ConvertDto2MatchActionModel(uriMatchAction.Action)
	matchResultModel, err := m.matchUsecase.MatchAction(ctx, *accountID, postMatchAction.TargetID, matchActionModel)
	if err != nil {
		log.Error("fail to perform match action", logger.FError(err))
		failResponse(ctx, http.StatusInternalServerError, dto.FailProcessRequestError, err)
		return
	}
	matchResult := converter.ConvertModel2MatchActionResponse(matchResultModel)
	successResponse(ctx, http.StatusOK, matchResult)
}
