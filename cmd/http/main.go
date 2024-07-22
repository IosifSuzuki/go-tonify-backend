package main

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/route"
	"go-tonify-backend/internal/bootstrap"
	"go-tonify-backend/internal/repository"
	"go-tonify-backend/internal/usecase"
	"log"
)

func main() {
	app := bootstrap.App()
	r := gin.Default()
	profileRepository := repository.NewProfileRepository()
	profileUseCase := usecase.NewProfileUseCase(profileRepository, app.ContentTimeout)
	route.Setup(r, profileUseCase)

	if err := r.Run(app.Address()); err != nil {
		log.Fatalln(err)
	}
}
