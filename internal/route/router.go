package route

import (
	"go-mandelbrot/internal/handlers"
	"net/http"
)

func HandleMandelbrot(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/mandelbrot", handlers.GetMandelbrotImageHandler)
	// Questionable, but okay
	//				👇
	mux.HandleFunc("PUT /v1/mandelbrot/colors", handlers.RegenerateColorRegions)
	mux.Handle("GET /", http.FileServer(http.Dir("./static/pages")))
	mux.Handle("GET /static/scripts/", http.StripPrefix("/static/scripts/", http.FileServer(http.Dir("./static/scripts"))))
}
