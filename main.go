package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vaibsharma/ddia/handler"
)

func gracefulShutdown(server *http.Server, timeout time.Duration) {
	// This function handles graceful shutdown of the HTTP server.
	// It takes a pointer to an http.Server and a timeout duration as parameters.
	// The server is shut down gracefully within the specified timeout.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}
}

func startServer() {
	// This function starts the HTTP server.
	// It creates a new http.Server instance and sets up the routes.
	// The server listens on port 8080 and handles incoming requests.
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	http.HandleFunc("/", handler.HelloWorldHandler)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe(): %s\n", err)
		}
	}()

	fmt.Println("Server started on :8080")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	sig := <-c
	fmt.Printf("Received signal: %s\n", sig)

	gracefulShutdown(server, 5*time.Second)
}

func main() {
	// This is the main function of the program.
	// It serves as the entry point for the application.
	// You can add your code here to execute when the program runs.
	// For example, you can call other functions or perform operations.
	// This is a placeholder for your main logic.

	// Create a new HTTP server and start it.
	startServer()
}
