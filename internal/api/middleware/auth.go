package middleware

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/model"
	"go-tonify-backend/internal/service"
	"log"
	"net/http"
)

type Auth struct {
	container   container.Container
	authService service.AuthService
}

func NewAuth(container container.Container, authService service.AuthService) *Auth {
	return &Auth{
		container:   container,
		authService: authService,
	}
}

func (a *Auth) Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenText := ctx.GetHeader(model.AuthorizationHeaderKey)
		if len(tokenText) == 0 {
			err := model.EmptyAuthorizationHeaderKeyError
			log.Println(err)
			sendResponse(ctx, http.StatusUnauthorized, model.EmptyAuthorizationHeaderKeyError)
			return
		}
		accessClaimsToken, err := a.authService.ParseAccessAccountJWT(tokenText)
		if err != nil {
			log.Println(err)
			sendResponse(ctx, http.StatusUnauthorized, err)
			return
		}
		ctx.Set(model.AuthorizationTokenKey, accessClaimsToken)
		ctx.Next()
	}
}
