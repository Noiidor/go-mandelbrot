package route

import (
	"go-mandelbrot/internal/handlers"
	"net/http"
)

func HandleMandelbrot(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/mandelbrot", handlers.GetMandelbrotImageHandler)
}
