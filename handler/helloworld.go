package handler

import "net/http"

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	// This is a simple HTTP handler function that responds with "Hello, World!".
	// It takes an http.ResponseWriter and an http.Request as parameters.
	// The response is written to the ResponseWriter.
	w.Write([]byte("Hello, World!"))
}
