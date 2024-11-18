package bootstrap

import (
	"net"
	"os"
	"strconv"
	"time"
)

type Application struct {
	ServerAddr          string
	ServerPort          string
	SecureServerAddr    string
	SecureServerPort    string
	ContentTimeout      time.Duration // in sec
	DBConfig            DBConfig
	TelegramConfig      TelegramConfig
	S3Config            S3Config
	JWTSecretKey        string
	AccessJWTExpiresIn  time.Duration // in sec
	RefreshJWTExpiresIn time.Duration // in sec
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Mode     string
}

type TelegramConfig struct {
	BotToken string
}

type S3Config struct {
	AttachmentBucket string
}

func App() Application {
	app := Application{}
	app.ServerAddr = os.Getenv("SERVER_HOST")
	app.ServerPort = os.Getenv("SERVER_PORT")
	app.SecureServerAddr = os.Getenv("SECURE_SERVER_HOST")
	app.SecureServerPort = os.Getenv("SECURE_SERVER_PORT")
	if contextTimeout, err := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT")); err == nil {
		app.ContentTimeout = time.Duration(contextTimeout) * time.Second
	} else {
		app.ContentTimeout = 3 * time.Second
	}
	app.JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
	if accessJWTExpiresIn, err := strconv.Atoi(os.Getenv("ACCESS_JWT_EXPIRES_IN")); err == nil {
		app.AccessJWTExpiresIn = time.Duration(accessJWTExpiresIn) * time.Second
	} else {
		app.AccessJWTExpiresIn = 1440 * time.Second
	}
	if refreshJWTExpiresIn, err := strconv.Atoi(os.Getenv("REFRESH_JWT_EXPIRES_IN")); err == nil {
		app.RefreshJWTExpiresIn = time.Duration(refreshJWTExpiresIn) * time.Second
	} else {
		app.RefreshJWTExpiresIn = 4320 * time.Second
	}

	app.DBConfig.Host = os.Getenv("DB_HOST")
	if dbPort, err := strconv.Atoi(os.Getenv("DB_PORT")); err == nil {
		app.DBConfig.Port = dbPort
	} else {
		app.DBConfig.Port = 5432
	}
	app.DBConfig.User = os.Getenv("DB_USER")
	app.DBConfig.Password = os.Getenv("DB_PASSWORD")
	app.DBConfig.Name = os.Getenv("DB_NAME")
	app.DBConfig.Mode = os.Getenv("DB_MODE")

	app.TelegramConfig.BotToken = os.Getenv("TELEGRAM_BOT_TOKEN")

	app.S3Config.AttachmentBucket = os.Getenv("S3_ATTACHMENT_BUCKET")

	return app
}

func (s *Application) Address() string {
	return net.JoinHostPort(s.ServerAddr, s.ServerPort)
}

func (s *Application) SecureAddress() string {
	return net.JoinHostPort(s.SecureServerAddr, s.SecureServerPort)
}
