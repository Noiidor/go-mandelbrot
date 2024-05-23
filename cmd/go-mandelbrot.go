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
		Addr:              "localhost:5050",
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	server.ListenAndServe()
}
