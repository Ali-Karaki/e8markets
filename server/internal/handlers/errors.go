package handlers

import (
	"errors"
	"net/http"

	"github.com/Ali-Karaki/e8markets/server/internal/apperr"
	"github.com/Ali-Karaki/e8markets/server/internal/httpx"
	"github.com/Ali-Karaki/e8markets/server/internal/tradelocker"
)

func writeAppError(w http.ResponseWriter, err *apperr.Error) {
	httpx.WriteAppError(w, err)
}

func writeValidationError(w http.ResponseWriter, message string) {
	writeAppError(w, apperr.Validation(message))
}

func writeInternalError(w http.ResponseWriter, message string) {
	writeAppError(w, apperr.Internal(message))
}

func writeTradeLockerError(w http.ResponseWriter, err error) {
	var tlErr *tradelocker.Error
	if errors.As(err, &tlErr) {
		writeAppError(w, mapTradeLockerError(tlErr))
		return
	}

	writeAppError(w, apperr.UpstreamError())
}

func mapTradeLockerError(tlErr *tradelocker.Error) *apperr.Error {
	switch tlErr.Code {
	case apperr.CodeUpstreamUnavailable:
		return apperr.UpstreamUnavailable()
	case apperr.CodeUpstreamRateLimited:
		return apperr.UpstreamRateLimited()
	case apperr.CodeUpstreamMalformed:
		return apperr.UpstreamMalformed()
	default:
		return apperr.UpstreamError()
	}
}
