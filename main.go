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

var (
	stop chan os.Signal
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

func SimulateCrash(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Simulating crash..."))
	stop <- os.Interrupt
}

func setupRoutes() *http.ServeMux {
	// This function sets up the HTTP routes for the server.
	// It returns a pointer to an http.ServeMux, which is a request multiplexer.
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// This is a health check endpoint.
		// It responds with a 200 OK status to indicate that the server is healthy.
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/", handler.HelloWorldHandler)
	mux.HandleFunc("/crash", SimulateCrash)
	return mux
}

func startServer() {
	// This function starts the HTTP server.
	// It creates a new http.Server instance and sets up the routes.
	// The server listens on port 8080 and handles incoming requests.

	// Server routes
	mux := setupRoutes()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe(): %s\n", err)
		}
	}()

	fmt.Println("Server started on :8080")

	stop = make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	sig := <-stop
	fmt.Printf("Received signal: %s\n", sig)

	gracefulShutdown(server, 5*time.Second)
}

func main() {
	// This is the main function of the program.
	// It serves as the entry point for the application.
	// You can add your code here to execute when the program runs.
	// For example, you can call other functions or perform operations.
	// This is a placeholder for your main logic.

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	// Create a new HTTP server and start it.
	startServer()

}
