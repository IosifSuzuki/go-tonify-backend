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
	return app
}

func (s *Application) Address() string {
	return net.JoinHostPort(s.ServerAddr, s.ServerPort)
}
