package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-mandelbrot/pkg/messages"
	"go-mandelbrot/pkg/service"
	"go-mandelbrot/pkg/templates"
	"net/http"
	"strconv"
)

func GetMandelbrotImageHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	pointXRaw := queryParams.Get("point_x")
	pointX, err := strconv.ParseFloat(pointXRaw, 64)
	if err != nil {
		errorJsonResponse(w, err, http.StatusBadRequest, "Cannot parse query parameter pointX.")
		return
	}
	pointYRaw := queryParams.Get("point_y")
	pointY, err := strconv.ParseFloat(pointYRaw, 64)
	if err != nil {
		errorJsonResponse(w, err, http.StatusBadRequest, "Cannot parse query parameter pointY.")
		return
	}
	zoomRaw := queryParams.Get("zoom")
	zoom, err := strconv.ParseUint(zoomRaw, 10, 64)
	if err != nil {
		errorJsonResponse(w, err, http.StatusBadRequest, "Cannot parse query parameter zoom.")
		return
	}
	resolutionWidthRaw := queryParams.Get("resolution_width")
	resolutionWidth, err := strconv.ParseUint(resolutionWidthRaw, 10, 32)
	if err != nil {
		errorJsonResponse(w, err, http.StatusBadRequest, "Cannot parse query parameter resolutionWidth.")
		return
	}
	resolutionHeightRaw := queryParams.Get("resolution_height")
	resolutionHeight, err := strconv.ParseUint(resolutionHeightRaw, 10, 32)
	if err != nil {
		errorJsonResponse(w, err, http.StatusBadRequest, "Cannot parse query parameter resolutionHeight.")
		return
	}

	service.GenerateMandelbrotImage(pointX, pointY, zoom, uint32(resolutionWidth), uint32(resolutionHeight))
}

var templater = templates.NewTemplates("pages/*.html")

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	err := templater.Render(&buf, "index.html", nil, r.Context())
	if err != nil {
		errorJsonResponse(w, err, http.StatusInternalServerError, "Failed to render index.html")
		return
	}

	buf.WriteTo(w)
}

func errorJsonResponse(w http.ResponseWriter, err error, code int, message string) {
	response := messages.ErrorsResponse{
		Errors: make([]messages.ErrorMessage, 1),
	}

	errorMessage := messages.ErrorMessage{
		Error:   fmt.Sprint(code),
		Message: message,
		Detail:  err.Error(),
	}

	response.Errors = append(response.Errors, errorMessage)

	w.WriteHeader(code)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	}
}
