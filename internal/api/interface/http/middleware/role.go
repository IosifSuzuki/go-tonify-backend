package middleware

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/api/interface/http/v1/converter"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain/account/usecase"
	"go-tonify-backend/pkg/logger"
	"net/http"
)

type Role struct {
	container      container.Container
	accountUsecase usecase.Account
}

func NewRole(container container.Container, accountUsecase usecase.Account) *Role {
	return &Role{
		container:      container,
		accountUsecase: accountUsecase,
	}
}

func (r *Role) Authorization(role dto.Role) gin.HandlerFunc {
	log := r.container.GetLogger()
	return func(ctx *gin.Context) {
		accountID, err := getAccountID(ctx)
		if err != nil {
			log.Error("fail to get account id", logger.FError(err))
			abortWithResponse(ctx, http.StatusUnauthorized, err)
			return
		}
		role := converter.ConvertDto2RoleModel(role)
		hasExpectedRole, err := r.accountUsecase.AccountHasRole(ctx, *accountID, role)
		if err != nil {
			log.Error(
				"fail to check role of account by id",
				logger.FError(err),
				logger.F("account_id", *accountID),
			)
			abortWithResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		if !hasExpectedRole {
			log.Error("expected other role", logger.F("expected_role", role))
			abortWithResponse(ctx, http.StatusForbidden, dto.RoleExpectedError)
			return
		}
		ctx.Next()
	}
}
