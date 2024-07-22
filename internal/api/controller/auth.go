package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
}

func (a *AuthController) SignIn(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func (a *AuthController) SignUp(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
