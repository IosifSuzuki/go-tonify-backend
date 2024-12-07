package container

import (
	"database/sql"
	"go-tonify-backend/internal/infrastructure/config"
	"go-tonify-backend/pkg/logger"
	"time"
)

type Container interface {
	GetLogger() logger.Logger
	GetTelegramBotToken() string
	GetAWSConfig() *config.AWS
	GetDBConnection() *sql.DB
	GetJWTSecretKey() string
	GetServerConfig() *config.Server
	GetAccessJWTExpiresIn() time.Duration
	GetRefreshJWTExpiresIn() time.Duration
}

type container struct {
	config *config.Config
	conn   *sql.DB
	logger logger.Logger
}

func NewContainer(
	logger logger.Logger,
	config *config.Config,
	conn *sql.DB,
) Container {
	return &container{
		config: config,
		conn:   conn,
		logger: logger,
	}
}

func (c *container) GetTelegramBotToken() string {
	return c.config.Telegram.BotToken
}

func (c *container) GetAWSConfig() *config.AWS {
	return c.config.AWS
}

func (c *container) GetDBConnection() *sql.DB {
	return c.conn
}

func (c *container) GetJWTSecretKey() string {
	return c.config.Server.JWTSecretKey
}

func (c *container) GetAccessJWTExpiresIn() time.Duration {
	return c.config.Server.AccessJWTExpiresIn
}

func (c *container) GetRefreshJWTExpiresIn() time.Duration {
	return c.config.Server.RefreshJWTExpiresIn
}

func (c *container) GetServerConfig() *config.Server {
	return c.config.Server
}

func (c *container) GetLogger() logger.Logger {
	return c.logger
}
