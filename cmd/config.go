package main

import (
	"context"
	"log"
	"sync"

	"github.com/imjenal/transaction-service/config"
	"github.com/imjenal/transaction-service/pkg/validator"
	"github.com/spf13/viper"
)

// The const block defines the keys for the key name in the .env file
const (
	keyAddress = "ADDR"
	keyPort    = "SERVER_PORT"
	keyEnv     = "ENVIRONMENT"

	keyDBHost     = "DB_HOST"
	keyDBPort     = "DB_PORT"
	keyDBUser     = "DB_USER"
	keyDBPassword = "DB_PASSWORD"
	keyDBName     = "DB_NAME"
)

// App Stores all the app config. The config is read from the .env file present in the project root.
type App struct {
	Server   *config.Server `validate:"required"`
	Database *config.DB     `validate:"required"`
}

var (
	configs *App
	once    sync.Once // the once is used to build a singleton pattern i.e, the data is only read once from .env
)

const envFileName = ".env"

// GetConfig reads the .env file and fetches the App config. Only the first function call reads the data from .env,
// subsequent calls just return the previously read data.
func GetConfig() *App {
	once.Do(func() {
		config.Read(envFileName, keyEnv)
		configs = &App{
			Server: &config.Server{
				Port:        viper.GetInt(keyPort),
				Address:     viper.GetString(keyAddress),
				Environment: config.Environment(viper.GetString(keyEnv)),
			},
			Database: &config.DB{
				Host:     viper.GetString(keyDBHost),
				Port:     viper.GetString(keyDBPort),
				User:     viper.GetString(keyDBUser),
				Password: viper.GetString(keyDBPassword),
				Name:     viper.GetString(keyDBName),
			},
		}

		validatr := validator.New()

		// Validate the struct to ensure all values are present and are correct
		result, err := validatr.IsValidStruct(context.Background(), configs)
		if err != nil {
			log.Fatalf("Error validation the config: %v", err)
		}

		if !result.Valid {
			log.Fatalf("Required config variables are missing or the value is invalid: %validatr", result.Fields)
		}
	})

	return configs
}
