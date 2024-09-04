package main

import (
	"context"
	"log"

	"github.com/imjenal/transaction-service/internal/app"

	"github.com/imjenal/transaction-service/api"
	"github.com/imjenal/transaction-service/internal/db"
	"github.com/imjenal/transaction-service/internal/server"
	"github.com/imjenal/transaction-service/pkg/http/request"
	"github.com/imjenal/transaction-service/pkg/http/response"
	"github.com/imjenal/transaction-service/pkg/validator"
)

var (
	Version = "0.0.0" // Version stores the binary version of the Git tag. It is populated using LDFlags
	Name    = "pismo"
)

func main() {
	// Set the version of the application.It will be used whenever call to app.Version() is made
	// This is the first thing that is done in the main function so that the version is set before any other function is called
	app.SetVersion(Version)

	config := GetConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize DB
	conn, err := db.GetConnection(ctx, &db.Config{
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		User:     config.Database.User,
		Password: config.Database.Password,
		Name:     config.Database.Name,
		Migrate:  true, // Always migrate the database
		SeedDB:   string(config.Server.Environment) == "DEV",
	})

	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return
	}

	// Defer closing the database connection, so that it is closed when the main function exits
	defer conn.Conn.Close()

	jsonWriter := response.NewJSONWriter()
	v := validator.New()

	// The params struct is passed to server, it is used when creating all the routes/handlers
	// It contains most of the top level dependencies
	// We initialize it here, so that we can have a single source of truth for all the dependencies
	// and also make is easier to test the code by passing in a mock implementation of the dependencies
	// instead of the actual implementation
	params := &api.Params{
		DB:        conn,
		Reader:    request.NewReader(jsonWriter, v),
		Writer:    jsonWriter,
		Validator: v,
	}

	serverConfig := &server.Config{
		Port: config.Server.Port,
		Host: config.Server.Address,
	}

	s := server.New(serverConfig, params) // Initialize the server
	s.Listen()                            // Start the server, this starts listening for incoming request

	// Wait for the server to shut down, this is a blocking call
	// It happens when the server receives a signal to shut down(eg: SIGINT) control is returned to the main function only after the server is shut down
	s.WaitForShutdown()
}
