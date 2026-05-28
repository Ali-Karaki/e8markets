package apperr

import "net/http"

type Error struct {
	Code       string
	Message    string
	HTTPStatus int
}

func (e *Error) Error() string {
	return e.Message
}

func newError(code, message string, status int) *Error {
	return &Error{Code: code, Message: message, HTTPStatus: status}
}

func InvalidCredentials() *Error {
	return newError(CodeInvalidCredentials, "Invalid credentials", http.StatusUnauthorized)
}

func NotAuthenticated() *Error {
	return newError(CodeNotAuthenticated, "Not authenticated", http.StatusUnauthorized)
}

func InvalidSession() *Error {
	return newError(CodeInvalidSession, "Invalid session", http.StatusUnauthorized)
}

func SessionExpired() *Error {
	return newError(CodeSessionExpired, "Session expired", http.StatusUnauthorized)
}

func Validation(message string) *Error {
	return newError(CodeValidationError, message, http.StatusBadRequest)
}

func Internal(message string) *Error {
	return newError(CodeInternalError, message, http.StatusInternalServerError)
}

func UpstreamUnavailable() *Error {
	return newError(CodeUpstreamUnavailable, "TradeLocker is temporarily unavailable", http.StatusServiceUnavailable)
}

func UpstreamRateLimited() *Error {
	return newError(CodeUpstreamRateLimited, "Too many requests. Try again shortly", http.StatusTooManyRequests)
}

func UpstreamError() *Error {
	return newError(CodeUpstreamError, "Could not complete the request with TradeLocker", http.StatusBadGateway)
}

func UpstreamMalformed() *Error {
	return newError(CodeUpstreamMalformed, "Received an unexpected response from TradeLocker", http.StatusBadGateway)
}
