package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/model"
	"go-tonify-backend/internal/service"
	"net/http"
)

type AuthController struct {
	container   container.Container
	AuthService service.AuthService
}

// AccountSignUp godoc
//
//	@Summary		account sign up
//	@Description	record account to db and return pairs jwt tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.CreateAccount	true	"account payload"
//	@Success		201		{object}	model.PairToken		"pair token"
//	@Failure		400		"bad parameters"
//	@Failure		500		"internal error"
//	@Router			/auth/account/sign-up [post]
func (a *AuthController) AccountSignUp(ctx *gin.Context) {
	var createAccount model.CreateAccount
	if err := ctx.ShouldBindJSON(&createAccount); err != nil {
		sendError(ctx, model.ParametersBadRequestError, http.StatusBadRequest)
		return
	}
	accountID, err := a.AuthService.CreateAccount(context.Background(), &createAccount)
	if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	pairToken, err := a.AuthService.GenerateAccountJWT(context.Background(), *accountID)
	if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	sendResponseWithStatus(ctx, pairToken, http.StatusCreated)
}
