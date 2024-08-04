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
