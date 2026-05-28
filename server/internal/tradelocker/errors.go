package tradelocker

import "github.com/Ali-Karaki/e8markets/server/internal/apperr"

type Error struct {
	Code           string
	UpstreamStatus int
	Message        string
}

func (e *Error) Error() string {
	return e.Message
}

func newTradeLockerError(code string, upstreamStatus int, message string) *Error {
	return &Error{Code: code, UpstreamStatus: upstreamStatus, Message: message}
}

func upstreamUnavailable(message string) *Error {
	return newTradeLockerError(apperr.CodeUpstreamUnavailable, 0, message)
}

func upstreamRateLimited(upstreamStatus int) *Error {
	return newTradeLockerError(apperr.CodeUpstreamRateLimited, upstreamStatus, "upstream rate limited")
}

func upstreamError(upstreamStatus int) *Error {
	return newTradeLockerError(apperr.CodeUpstreamError, upstreamStatus, "upstream error")
}

func upstreamMalformed(message string) *Error {
	return newTradeLockerError(apperr.CodeUpstreamMalformed, 0, message)
}
