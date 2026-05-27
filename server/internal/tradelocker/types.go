package tradelocker

import (
	"encoding/json"
	"strconv"
	"time"
)

// flexFloat64 unmarshals JSON numbers or numeric strings from TradeLocker.
// OpenAPI declares several numeric fields as "number", but the live API often
// returns them as strings (e.g. accountBalance "100000.00").
type flexFloat64 float64

func (f *flexFloat64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	var n float64
	if err := json.Unmarshal(data, &n); err == nil {
		*f = flexFloat64(n)
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*f = flexFloat64(v)
	return nil
}

func (f flexFloat64) Float64() float64 {
	return float64(f)
}

// API-mapped types below mirror TradeLocker Public API schemas.
// Do not rename json tags or fields unless the upstream API changes.
// Docs: https://public-api.tradelocker.com/docs/getting-started

// LoginRequest is the body for POST /auth/jwt/token.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Server   string `json:"server"`
}

// TokenResponse is returned by POST /auth/jwt/token and POST /auth/jwt/refresh.
type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpireDate   string `json:"expireDate"`
}

// RefreshRequest is the body for POST /auth/jwt/refresh.
type RefreshRequest struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// Account is an entry in GET /auth/jwt/all-accounts.
// OpenAPI: https://public-api.tradelocker.com/reference/getallaccounts.md
type Account struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
	AccNum   string `json:"accNum"`
	// accountBalance is returned by the live API; aaccountBalance is the OpenAPI typo field.
	AccountBalance  flexFloat64 `json:"accountBalance,omitempty"`
	AAccountBalance flexFloat64 `json:"aaccountBalance,omitempty"`
}

// AccountsResponse is the top-level shape of GET /auth/jwt/all-accounts.
type AccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

// TLResponse is the { "d": ..., "s": "ok" } envelope used by /trade/* endpoints.
type TLResponse[T any] struct {
	D T      `json:"d"`
	S string `json:"s"`
}

// AccountsData is the "d" payload when accounts are nested in a trade-style response.
type AccountsData struct {
	Accounts []Account `json:"accounts"`
}

// AccountStateData is the "d" payload for GET /trade/accounts/{accountId}/state.
// OpenAPI: https://public-api.tradelocker.com/reference/getstate.md
type AccountStateData struct {
	AccountDetailsData []flexFloat64 `json:"accountDetailsData"`
}

// ConfigData is the "d" payload for GET /trade/config.
// OpenAPI: https://public-api.tradelocker.com/reference/getconfigusingget.md
type ConfigData struct {
	AccountDetailsConfig *ColumnConfig `json:"accountDetailsConfig"`
}

// ColumnConfig maps accountDetailsData array indices to column id values (PanelConfig in OpenAPI).
type ColumnConfig struct {
	ID      string   `json:"id,omitempty"`
	Title   string   `json:"title,omitempty"`
	Columns []Column `json:"columns"`
}

type Column struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// Route is an instrument routing entry from GET /trade/accounts/{accountId}/instruments.
// OpenAPI: https://public-api.tradelocker.com/reference/getinstruments.md
type Route struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

// Instrument is an entry in GET /trade/accounts/{accountId}/instruments.
type Instrument struct {
	ID                   int64   `json:"id"`
	Name                 string  `json:"name"`
	Description          string  `json:"description,omitempty"`
	Type                 string  `json:"type"`
	TradableInstrumentID int64   `json:"tradableInstrumentId"`
	TradingExchange      string  `json:"tradingExchange"`
	MarketDataExchange   string  `json:"marketDataExchange"`
	Routes               []Route `json:"routes,omitempty"`
}

// InstrumentsData is the "d" payload for GET /trade/accounts/{accountId}/instruments.
type InstrumentsData struct {
	Instruments []Instrument `json:"instruments"`
}

// InstrumentSummary is a frontend-friendly instrument shape returned by our API.
type InstrumentSummary struct {
	ID                   int64  `json:"id"`
	TradableInstrumentID int64  `json:"tradableInstrumentId"`
	Name                 string `json:"name"`
	Description          string `json:"description,omitempty"`
	Type                 string `json:"type"`
	TradingExchange      string `json:"tradingExchange"`
	MarketDataExchange   string `json:"marketDataExchange"`
	TradeRouteID         int64  `json:"tradeRouteId,omitempty"`
	InfoRouteID          int64  `json:"infoRouteId,omitempty"`
}

// Tokens is an internal app type (parsed ExpireDate); not a raw TL response shape.
type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
