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
//	@Router			/auth/sign-up [post]
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

// AccountSignIn godoc
//
//	@Summary		account sign in
//	@Description	process of authorization to system through provided credentials
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.Credential	true	"credential"
//	@Success		200		{object}	model.PairToken		"pair token"
//	@Failure		400		"bad parameters"
//	@Failure		401		"incorrect or missing credentials"
//	@Failure		500		"internal error"
//	@Router			/auth/sign-in [post]
func (a *AuthController) AccountSignIn(ctx *gin.Context) {
	var credential model.Credential
	if err := ctx.ShouldBindJSON(&credential); err != nil {
		sendError(ctx, model.ParametersBadRequestError, http.StatusBadRequest)
		return
	}
	account, err := a.AuthService.AuthorizationAccount(ctx, &credential)
	if err == model.AccountNotExistsError {
		sendError(ctx, err, http.StatusUnauthorized)
		return
	} else if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	pairToken, err := a.AuthService.GenerateAccountJWT(context.Background(), *account.ID)
	if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	sendResponseWithStatus(ctx, pairToken, http.StatusOK)
}
