package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/Ali-Karaki/e8markets/server/internal/apperr"
)

type errorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func CodedError(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, errorResponse{Error: message, Code: code})
}

func WriteAppError(w http.ResponseWriter, err *apperr.Error) {
	CodedError(w, err.HTTPStatus, err.Code, err.Message)
}
