package main

import (
	"go-mandelbrot/internal/route"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	route.HandleMandelbrot(mux)

	server := http.Server{
		Addr:              "localhost:5050", //TODO: change to config or flag
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	server.ListenAndServe()
}
