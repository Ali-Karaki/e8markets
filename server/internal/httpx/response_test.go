package httpx

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestCodedError_includesCode(t *testing.T) {
	rec := httptest.NewRecorder()
	CodedError(rec, 502, "upstream_error", "Could not complete the request with TradeLocker")

	if rec.Code != 502 {
		t.Fatalf("status = %d", rec.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}

	if body["error"] != "Could not complete the request with TradeLocker" {
		t.Fatalf("error = %q", body["error"])
	}
	if body["code"] != "upstream_error" {
		t.Fatalf("code = %q", body["code"])
	}
}
