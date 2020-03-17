package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Respond converts a Go value to JSON and sends it to the client.
func (h *Handler) respond(w http.ResponseWriter, data interface{}, statusCode int) {
	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent || data == nil {
		w.WriteHeader(statusCode)
		return
	}

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if err := json.NewEncoder(w).Encode(&data); err != nil {
		h.respondError(w, err)
	}

	return
}

func (h *Handler) respondError(w http.ResponseWriter, err error) {
	er := ErrorResponse{
		Message:       err.Error(),
		StatusCode:    http.StatusInternalServerError,
		StatusMessage: http.StatusText(http.StatusInternalServerError),
	}

	if err, ok := err.(APIError); ok {
		er.Message = err.Message
		er.StatusCode = err.StatusCode
		er.StatusMessage = http.StatusText(err.StatusCode)
	}

	h.log.WithFields(logrus.Fields{
		"code": er.StatusCode,
	}).Error(er.Message)

	h.respond(w, er, er.StatusCode)
}
