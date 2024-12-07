package config

import (
	"go-tonify-backend/internal/domain/entity"
	"os"
	"strconv"
	"sync"
)

type PostgreSQL struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Mode     string
}

var (
	postgreSQLInstance *PostgreSQL
	postgreSQLErr      error
	postgreSQLOnce     sync.Once
)

func GetPostgreSQL() (*PostgreSQL, error) {
	postgreSQLOnce.Do(func() {
		var (
			instance PostgreSQL
			ok       bool
			err      error
		)
		instance.Host, ok = os.LookupEnv("POSTGRESQL_HOST")
		if !ok {
			postgreSQLErr = entity.NilError
			return
		}
		portText, ok := os.LookupEnv("POSTGRESQL_PORT")
		if !ok {
			postgreSQLErr = entity.NilError
			return
		}
		instance.Port, err = strconv.Atoi(portText)
		if err != nil {
			postgreSQLErr = entity.ConvertStringToIntError
			return
		}
		instance.User, ok = os.LookupEnv("POSTGRESQL_USER")
		if !ok {
			postgreSQLErr = entity.NilError
			return
		}
		instance.Password, ok = os.LookupEnv("POSTGRESQL_PASSWORD")
		if !ok {
			postgreSQLErr = entity.NilError
			return
		}
		instance.Name, ok = os.LookupEnv("POSTGRESQL_NAME")
		if !ok {
			postgreSQLErr = entity.NilError
			return
		}
		instance.Mode, ok = os.LookupEnv("POSTGRESQL_MODE")
		if !ok {
			postgreSQLErr = entity.NilError
			return
		}
		postgreSQLInstance = &instance
	})
	return postgreSQLInstance, postgreSQLErr
}
