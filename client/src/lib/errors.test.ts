import { describe, expect, it } from "vitest";

import { ApiError } from "./api";
import { ApiErrorCode, apiErrorMessage, isSessionExpiredCode } from "./errors";

describe("apiErrorMessage", () => {
  it("maps known error codes", () => {
    expect(
      apiErrorMessage(
        new ApiError(503, "", ApiErrorCode.UpstreamUnavailable),
        "fallback",
      ),
    ).toBe("TradeLocker is temporarily unavailable. Try again later.");
  });

  it("uses ApiError message when code is not mapped", () => {
    expect(
      apiErrorMessage(
        new ApiError(400, "Bad input", ApiErrorCode.ValidationError),
        "fallback",
      ),
    ).toBe("Bad input");
  });

  it("falls back for unknown errors", () => {
    expect(apiErrorMessage(new Error("boom"), "fallback")).toBe("fallback");
  });
});

describe("isSessionExpiredCode", () => {
  it("detects session-related codes", () => {
    expect(isSessionExpiredCode(ApiErrorCode.SessionExpired)).toBe(true);
    expect(isSessionExpiredCode(ApiErrorCode.InvalidSession)).toBe(true);
    expect(isSessionExpiredCode(ApiErrorCode.NotAuthenticated)).toBe(true);
    expect(isSessionExpiredCode(ApiErrorCode.InvalidCredentials)).toBe(false);
  });
});
