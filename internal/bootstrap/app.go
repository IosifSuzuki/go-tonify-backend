package bootstrap

import (
	"net"
	"os"
	"strconv"
	"time"
)

type Application struct {
	ServerAddr     string
	ServerPort     string
	ContentTimeout time.Duration // in sec
	DBConfig       DBConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Mode     string
}

func App() Application {
	app := Application{}
	app.ServerAddr = os.Getenv("SERVER_HOST")
	app.ServerPort = os.Getenv("SERVER_PORT")
	if contextTimeout, err := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT")); err == nil {
		app.ContentTimeout = time.Duration(contextTimeout)
	} else {
		app.ContentTimeout = 3
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
	return app
}

func (s *Application) Address() string {
	return net.JoinHostPort(s.ServerAddr, s.ServerPort)
}
