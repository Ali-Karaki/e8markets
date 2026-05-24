package tradelocker

import "time"

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
type Account struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Currency        string  `json:"currency"`
	Status          string  `json:"status"`
	AccNum          string  `json:"accNum"`
	// @TODO:Ali check if we can drop AccountBalance
	AccountBalance  float64 `json:"accountBalance,omitempty"`
	AAccountBalance float64 `json:"aaccountBalance,omitempty"` // typo in TL OpenAPI: aaccountBalance
}

// AccountsResponse is the top-level shape of GET /auth/jwt/all-accounts.
type AccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

// TLResponse is the { "d": ..., "s": "ok" } envelope used by /trade/* endpoints.
type TLResponse[T interface{}] struct {
	D T      `json:"d"`
	S string `json:"s"`
}

// AccountsData is the "d" payload when accounts are nested in a trade-style response.
type AccountsData struct {
	Accounts []Account `json:"accounts"`
}

// AccountStateData is the "d" payload for GET /trade/accounts/{accountId}/state.
type AccountStateData struct {
	AccountDetailsData []float64 `json:"accountDetailsData"`
}

// ConfigData is the "d" payload for GET /trade/config.
type ConfigData struct {
	AccountDetailsConfig *ColumnConfig `json:"accountDetailsConfig"`
}

// ColumnConfig and Column map accountDetailsData array indices to column id values.
type ColumnConfig struct {
	Columns []Column `json:"columns"`
}

type Column struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// Tokens is an internal app type (parsed ExpireDate); not a raw TL response shape.
type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
