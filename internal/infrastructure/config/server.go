package config

import (
	"go-tonify-backend/internal/domain/entity"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	Addr                string
	Port                string
	SecureAddr          string
	SecurePort          string
	JWTSecretKey        string
	LoggerLevel         int
	AccessJWTExpiresIn  time.Duration // in sec
	RefreshJWTExpiresIn time.Duration // in sec
}

var (
	serverInstance *Server
	serverErr      error
	serverOnce     sync.Once
)

func GetServer() (*Server, error) {
	serverOnce.Do(func() {
		var (
			ok       bool
			instance Server
			err      error
		)
		instance.Addr, ok = os.LookupEnv("SERVER_HOST")
		if !ok {
			serverErr = entity.NilError
			return
		}
		instance.Port, ok = os.LookupEnv("SERVER_PORT")
		if !ok {
			serverErr = entity.NilError
			return
		}
		instance.SecureAddr, ok = os.LookupEnv("SECURE_SERVER_HOST")
		if !ok {
			serverErr = entity.NilError
			return
		}
		instance.SecurePort, ok = os.LookupEnv("SECURE_SERVER_PORT")
		if !ok {
			serverErr = entity.NilError
			return
		}
		loggerLevelText, ok := os.LookupEnv("SERVER_LOGGER_LEVEL")
		if !ok {
			serverErr = entity.NilError
			return
		}
		instance.LoggerLevel, err = strconv.Atoi(loggerLevelText)
		if err != nil {
			serverErr = err
			return
		}
		accessJWTExpiresInText, ok := os.LookupEnv("ACCESS_JWT_EXPIRES_IN")
		if !ok {
			serverErr = entity.NilError
			return
		}
		accessJWTExpiresIn, err := strconv.Atoi(accessJWTExpiresInText)
		if err != nil {
			serverErr = entity.ConvertStringToIntError
			return
		}
		instance.AccessJWTExpiresIn = time.Duration(accessJWTExpiresIn) * time.Second
		refreshJWTExpiresInText, ok := os.LookupEnv("REFRESH_JWT_EXPIRES_IN")
		if !ok {
			serverErr = entity.NilError
			return
		}
		refreshJWTExpiresIn, err := strconv.Atoi(refreshJWTExpiresInText)
		if err != nil {
			serverErr = entity.ConvertStringToIntError
			return
		}
		instance.RefreshJWTExpiresIn = time.Duration(refreshJWTExpiresIn) * time.Second
		serverInstance = &instance
	})
	return serverInstance, serverErr
}

func (s *Server) Address() string {
	return net.JoinHostPort(s.Addr, s.Port)
}

func (s *Server) SecureAddress() string {
	return net.JoinHostPort(s.SecureAddr, s.SecurePort)
}
