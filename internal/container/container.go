package container

import (
	"database/sql"
	"go-tonify-backend/internal/bootstrap"
	"go-tonify-backend/pkg/logger"
	"time"
)

type Container interface {
	GetLogger() logger.Logger
	GetTelegramConfig() bootstrap.TelegramConfig
	GetContentTimeout() time.Duration
	GetDBConnection() *sql.DB
	GetJWTSecretKey() string
	GetAccessJWTExpiresIn() time.Duration
	GetRefreshJWTExpiresIn() time.Duration
}
type container struct {
	config bootstrap.Application
	conn   *sql.DB
	logger logger.Logger
}

func NewContainer(logger logger.Logger, config bootstrap.Application, conn *sql.DB) Container {
	return &container{
		config: config,
		conn:   conn,
		logger: logger,
	}
}

func (c *container) GetTelegramConfig() bootstrap.TelegramConfig {
	return c.config.TelegramConfig
}

func (c *container) GetContentTimeout() time.Duration {
	return c.config.ContentTimeout
}

func (c *container) GetDBConnection() *sql.DB {
	return c.conn
}

func (c *container) GetJWTSecretKey() string {
	return c.config.JWTSecretKey
}

func (c *container) GetAccessJWTExpiresIn() time.Duration {
	return c.config.AccessJWTExpiresIn
}

func (c *container) GetRefreshJWTExpiresIn() time.Duration {
	return c.config.RefreshJWTExpiresIn
}

func (c *container) GetLogger() logger.Logger {
	return c.logger
}
