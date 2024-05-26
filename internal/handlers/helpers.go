package handlers

import (
	"encoding/json"
	"fmt"
	"go-mandelbrot/pkg/messages"
	"net/http"
)

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
