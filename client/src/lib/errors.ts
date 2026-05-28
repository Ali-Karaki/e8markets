import { ApiError } from "./api";

export const ApiErrorCode = {
  InvalidCredentials: "invalid_credentials",
  NotAuthenticated: "not_authenticated",
  InvalidSession: "invalid_session",
  SessionExpired: "session_expired",
  ValidationError: "validation_error",
  UpstreamUnavailable: "upstream_unavailable",
  UpstreamRateLimited: "upstream_rate_limited",
  UpstreamError: "upstream_error",
  UpstreamMalformed: "upstream_malformed",
  InternalError: "internal_error",
} as const;

const codeMessages: Record<string, string> = {
  [ApiErrorCode.UpstreamUnavailable]:
    "TradeLocker is temporarily unavailable. Try again later.",
  [ApiErrorCode.UpstreamRateLimited]: "Too many requests. Try again shortly.",
  [ApiErrorCode.UpstreamError]:
    "Could not complete the request with TradeLocker.",
  [ApiErrorCode.UpstreamMalformed]:
    "Received an unexpected response from TradeLocker.",
  [ApiErrorCode.SessionExpired]: "Your session expired. Please sign in again.",
  [ApiErrorCode.NotAuthenticated]: "Please sign in to continue.",
  [ApiErrorCode.InvalidSession]:
    "Your session is invalid. Please sign in again.",
  [ApiErrorCode.InternalError]: "Something went wrong. Please try again.",
};

export function apiErrorMessage(error: unknown, fallback: string): string {
  if (error instanceof ApiError) {
    if (error.code && codeMessages[error.code]) {
      return codeMessages[error.code];
    }
    if (error.message) {
      return error.message;
    }
  }

  return fallback;
}

export function isSessionExpiredCode(code: string | undefined): boolean {
  return (
    code === ApiErrorCode.SessionExpired ||
    code === ApiErrorCode.InvalidSession ||
    code === ApiErrorCode.NotAuthenticated
  );
}
