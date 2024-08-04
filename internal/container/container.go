package container

import (
	"database/sql"
	"go-tonify-backend/internal/bootstrap"
	"time"
)

type Container interface {
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
}

func NewContainer(config bootstrap.Application, conn *sql.DB) Container {
	return &container{
		config: config,
		conn:   conn,
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
