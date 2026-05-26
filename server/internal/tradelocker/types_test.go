package tradelocker

import (
	"encoding/json"
	"testing"
)

func TestAccountsResponse_liveAPI(t *testing.T) {
	const raw = `{"accounts":[{"id":"2200766","name":"E8#demo","currency":"USD","accNum":"1","accountBalance":"100000.00","status":"ACTIVE"}]}`

	var resp AccountsResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Accounts[0].AccountBalance.Float64() != 100000 {
		t.Fatalf("accountBalance = %v", resp.Accounts[0].AccountBalance.Float64())
	}
}

func TestAccountsResponse_openAPI_aaccountBalance(t *testing.T) {
	const raw = `{"accounts":[{"id":"7080","accNum":"1","aaccountBalance":"2024.75"}]}`

	var resp AccountsResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Accounts[0].AAccountBalance.Float64() != 2024.75 {
		t.Fatalf("aaccountBalance = %v", resp.Accounts[0].AAccountBalance.Float64())
	}
}

func TestAccountStateData_openAPI_numeric(t *testing.T) {
	const raw = `{"d":{"accountDetailsData":[392084.93,391415.43,389608.865,0,392084.93,0,389608.865,0,0,1806.965,1806.965,100,0,0,144,1806.965,389608.865,5725.99,5585.57,140.42,131.12,11,-669.5,-739.5,1,0]},"s":"ok"}`

	var resp TLResponse[AccountStateData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.D.AccountDetailsData) != 26 {
		t.Fatalf("len = %d", len(resp.D.AccountDetailsData))
	}
	if resp.D.AccountDetailsData[0].Float64() != 392084.93 {
		t.Fatalf("balance = %v", resp.D.AccountDetailsData[0].Float64())
	}
}

func TestAccountStateData_stringValues(t *testing.T) {
	const raw = `{"d":{"accountDetailsData":["100000.00","100000.00","100000.00","0","100000.00"]},"s":"ok"}`

	var resp TLResponse[AccountStateData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.D.AccountDetailsData[2].Float64() != 100000 {
		t.Fatalf("availableFunds = %v", resp.D.AccountDetailsData[2].Float64())
	}
}

func TestConfigData_accountDetails(t *testing.T) {
	const raw = `{"s":"ok","d":{"accountDetailsConfig":{"id":"accountDetails","title":"Account details","columns":[{"id":"balance","description":"Balance"},{"id":"projectedBalance","description":"Projected"},{"id":"availableFunds","description":"Available"},{"id":"blockedBalance","description":"Blocked"},{"id":"cashBalance","description":"Cash"}]}}}`

	var resp TLResponse[ConfigData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	cfg := resp.D.AccountDetailsConfig
	if cfg == nil || cfg.ID != "accountDetails" || len(cfg.Columns) != 5 {
		t.Fatalf("config = %+v", cfg)
	}
}

func TestTokenResponse(t *testing.T) {
	const raw = `{"accessToken":"tok","refreshToken":"ref","expireDate":"2023-05-26T16:59:53.000Z"}`

	var resp TokenResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.AccessToken != "tok" || resp.ExpireDate == "" {
		t.Fatalf("tokens = %+v", resp)
	}
}

func TestMapAccountDetails(t *testing.T) {
	columns := []Column{
		{ID: "balance"},
		{ID: "projectedBalance"},
		{ID: "availableFunds"},
	}
	data := []flexFloat64{100000, 100000, 100000}

	state := mapAccountDetails(columns, data)
	if state["balance"] != 100000 || state["availableFunds"] != 100000 {
		t.Fatalf("state = %+v", state)
	}
}

func TestFlexFloat64_rejectsNonNumericString(t *testing.T) {
	var f flexFloat64
	if err := json.Unmarshal([]byte(`"not-a-number"`), &f); err == nil {
		t.Fatal("expected error")
	}
}
