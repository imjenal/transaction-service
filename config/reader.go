package config

import (
	"log"

	"github.com/spf13/viper"
)

// Read reads the .env file or the environment variables from the OS.
// If the file is not present and no OS envs are defined, or any other error in read, the application panics.
// It uses the spf13/viper library to read the .env file and the environment variables.
func Read(envFile, verifyEnvKey string) {
	//Set the config file name & type
	viper.SetConfigName(envFile)
	viper.SetConfigType("env")

	// Set the path to look for the config file
	viper.AddConfigPath(".")

	// Try to read the config file
	err := viper.ReadInConfig()
	if err == nil {
		// The config file is present, and the values were read
		log.Println("Environment variables successfully read from the .env")
	}

	// Reading from config file complete
	// Override with any environment variables passed through CLI
	viper.AutomaticEnv()

	// Check if the config is read from the environment variables
	if viper.GetString(verifyEnvKey) != "" {
		// The config is read from the environment variables
		log.Println("Environment variables successfully read from the OS")
		return
	}

	// The config is not read from the environment variables or the .env file
	log.Panicf("error reading the config file or the environment variables: %v", err)
}
