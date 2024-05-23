package main

import (
	"go-mandelbrot/internal/route"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	route.HandleMandelbrot(mux)

	server := http.Server{
		Addr:    "localhost:5050",
		Handler: mux,
	}

	server.ListenAndServe()
}
