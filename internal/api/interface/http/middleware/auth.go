package middleware

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain/account/usecase"
	"go-tonify-backend/pkg/logger"
	"net/http"
)

type Auth struct {
	container      container.Container
	accountUsecase usecase.Account
}

func NewAuth(container container.Container, accountUsecase usecase.Account) *Auth {
	return &Auth{
		container:      container,
		accountUsecase: accountUsecase,
	}
}

func (a *Auth) Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := a.container.GetLogger()
		tokenText := ctx.GetHeader(dto.AuthorizationHeaderKey)
		if len(tokenText) == 0 {
			log.Error("missing authorization header value")
			abortWithResponse(ctx, http.StatusUnauthorized, dto.MissingAuthorizationTokenError)
			return
		}
		accountID, err := a.accountUsecase.ParseAccessToken(tokenText)
		if err != nil {
			log.Error("fail to parse / validate parse access token", logger.FError(err))
			abortWithResponse(ctx, http.StatusUnauthorized, dto.ParseValidateTokenError)
			return
		}
		if accountID == nil {
			log.Error("account id has nil value", logger.FError(err))
			abortWithResponse(ctx, http.StatusUnauthorized, dto.NilError)
			return
		}
		ctx.Set(dto.AccountIDKey, *accountID)
		ctx.Next()
	}
}
