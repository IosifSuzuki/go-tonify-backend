package config

import "sync"

type Config struct {
	Server     *Server
	AWS        *AWS
	PostgreSQL *PostgreSQL
	Telegram   *Telegram
}

var (
	configOnce     sync.Once
	configInstance *Config
	configError    error
)

func GetConfig() (*Config, error) {
	configOnce.Do(func() {
		var (
			instance Config
			err      error
		)
		instance.Server, err = GetServer()
		if err != nil {
			configError = err
			return
		}
		instance.AWS, err = GetAWS()
		if err != nil {
			configError = err
			return
		}
		instance.Telegram, err = GetTelegram()
		if err != nil {
			configError = err
			return
		}
		instance.PostgreSQL, err = GetPostgreSQL()
		if err != nil {
			configError = err
			return
		}
		configInstance = &instance
	})
	return configInstance, configError
}
