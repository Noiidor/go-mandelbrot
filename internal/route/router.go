package route

import (
	"go-mandelbrot/internal/handlers"
	"net/http"
)

func HandleMandelbrot(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/mandelbrot", handlers.GetMandelbrotImageHandler)
	mux.Handle("GET /", http.FileServer(http.Dir("./static/pages")))
	mux.Handle("GET /static/scripts/", http.StripPrefix("/static/scripts/", http.FileServer(http.Dir("./static/scripts"))))
}
