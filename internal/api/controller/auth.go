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

// ClientSignUp godoc
//
//	@Summary		client sign up
//	@Description	record client to db and return pairs jwt tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.CreateClient	true	"client payload"
//	@Success		201		{object}	model.PairToken		"pair token"
//	@Failure		400		"bad parameters"
//	@Failure		500		"internal error"
//	@Router			/auth/client/sign-up [post]
func (a *AuthController) ClientSignUp(ctx *gin.Context) {
	var createClient model.CreateClient
	if err := ctx.ShouldBindJSON(&createClient); err != nil {
		sendError(ctx, model.ParametersBadRequestError, http.StatusBadRequest)
		return
	}
	clientID, err := a.AuthService.CreateClient(context.Background(), &createClient)
	if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	pairToken, err := a.AuthService.GenerateClientJWT(context.Background(), *clientID)
	if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	sendResponseWithStatus(ctx, pairToken, http.StatusCreated)
}

// FreelancerSignUp godoc
//
//	@Summary		freelancer sign up
//	@Description	record freelancer to db and return pairs jwt tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.CreateFreelancer	true	"freelancer payload"
//	@Success		201		{object}	model.PairToken			"pair token"
//	@Failure		400		"bad parameters"
//	@Failure		500		"internal error"
//	@Router			/auth/freelancer/sign-up [post]
func (a *AuthController) FreelancerSignUp(ctx *gin.Context) {
	var createFreelancer model.CreateFreelancer
	if err := ctx.ShouldBindJSON(&createFreelancer); err != nil {
		sendError(ctx, model.ParametersBadRequestError, http.StatusBadRequest)
		return
	}
	freelancerID, err := a.AuthService.CreateFreelancer(context.Background(), &createFreelancer)
	if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	pairToken, err := a.AuthService.GenerateFreelancerJWT(context.Background(), *freelancerID)
	if err != nil {
		sendError(ctx, err, http.StatusInternalServerError)
		return
	}
	sendResponseWithStatus(ctx, pairToken, http.StatusCreated)
}
