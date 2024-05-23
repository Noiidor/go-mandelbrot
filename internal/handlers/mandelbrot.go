package handlers

import (
	"go-mandelbrot/pkg/service"
	"net/http"
)

func GetMandelbrotImageHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	pointX := queryParams.Get("point_x")
	pointY := queryParams.Get("point_y")
	zoom := queryParams.Get("zoom")
	resolutionWidth := queryParams.Get("resolution_width")
	resolutionHeight := queryParams.Get("resolution_height")

	service.GenerateMandelbrotImage(pointX, pointY, zoom, resolutionWidth, resolutionHeight)
}
