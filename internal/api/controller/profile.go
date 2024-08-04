package controller

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/service"
	"net/http"
)

type ProfileController struct {
	ProfileUseCase service.ProfileService
}

func (p *ProfileController) GetProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	profile, err := p.ProfileUseCase.GetProfileByID(ctx, id)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, profile)
}
