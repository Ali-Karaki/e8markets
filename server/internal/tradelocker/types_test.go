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

func TestInstrumentsData_openAPI(t *testing.T) {
	const raw = `{"d":{"instruments":[{"id":1,"name":"EURUSD","tradableInstrumentId":1001,"type":"FOREX","tradingExchange":"FX","marketDataExchange":"FX","description":"Euro vs US Dollar","routes":[{"id":42,"type":"TRADE"},{"id":43,"type":"INFO"}]}]},"s":"ok"}`

	var resp TLResponse[InstrumentsData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.D.Instruments) != 1 {
		t.Fatalf("len = %d", len(resp.D.Instruments))
	}
	inst := resp.D.Instruments[0]
	if inst.Name != "EURUSD" || inst.TradableInstrumentID != 1001 {
		t.Fatalf("instrument = %+v", inst)
	}
}

func TestExtractRouteIDs(t *testing.T) {
	routes := []Route{
		{ID: 10, Type: "INFO"},
		{ID: 20, Type: "TRADE"},
	}
	trade, info := extractRouteIDs(routes)
	if trade != 20 || info != 10 {
		t.Fatalf("trade=%d info=%d", trade, info)
	}
}

func TestConfigData_positionsConfig(t *testing.T) {
	const raw = `{"s":"ok","d":{"positionsConfig":{"id":"positions","title":"Positions","columns":[{"id":"id"},{"id":"tradableInstrumentId"},{"id":"side"},{"id":"qty"},{"id":"avgPrice"}]}}}`

	var resp TLResponse[ConfigData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	cfg := resp.D.PositionsConfig
	if cfg == nil || cfg.ID != "positions" || len(cfg.Columns) != 5 {
		t.Fatalf("config = %+v", cfg)
	}
}

func TestPositionsData_liveAPI(t *testing.T) {
	const raw = `{"s":"ok","d":{"positions":[["360287970189927922","6109","948735","buy","0.01","1.85784",null,null,"1779927208000","-0.10",null]]}}`

	var resp TLResponse[PositionsData]
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.D.Positions) != 1 || len(resp.D.Positions[0]) != 11 {
		t.Fatalf("positions = %+v", resp.D.Positions)
	}
}

func TestMapPositions(t *testing.T) {
	columns := []Column{
		{ID: "id"},
		{ID: "tradableInstrumentId"},
		{ID: "routeId"},
		{ID: "side"},
		{ID: "qty"},
		{ID: "avgPrice"},
	}
	row := []json.RawMessage{
		[]byte(`1001`),
		[]byte(`2001`),
		[]byte(`42`),
		[]byte(`"buy"`),
		[]byte(`"1.5"`),
		[]byte(`1.08432`),
	}

	positions := mapPositions(columns, [][]json.RawMessage{row})
	if len(positions) != 1 {
		t.Fatalf("len = %d", len(positions))
	}
	p := positions[0]
	if p.ID != "1001" || p.TradableInstrumentID != 2001 || p.Side != "buy" || p.Qty != 1.5 {
		t.Fatalf("position = %+v", p)
	}
	if p.AvgPrice != 1.08432 {
		t.Fatalf("avgPrice = %v", p.AvgPrice)
	}
}

func TestMapPositions_liveRow(t *testing.T) {
	columns := []Column{
		{ID: "id"},
		{ID: "tradableInstrumentId"},
		{ID: "routeId"},
		{ID: "side"},
		{ID: "qty"},
		{ID: "avgPrice"},
		{ID: "stopLossId"},
		{ID: "takeProfitId"},
		{ID: "openDate"},
		{ID: "unrealizedPl"},
		{ID: "strategyId"},
	}
	row := []json.RawMessage{
		[]byte(`"360287970189927922"`),
		[]byte(`"6109"`),
		[]byte(`"948735"`),
		[]byte(`"buy"`),
		[]byte(`"0.01"`),
		[]byte(`"1.85784"`),
		[]byte(`null`),
		[]byte(`null`),
		[]byte(`"1779927208000"`),
		[]byte(`"-0.10"`),
		[]byte(`null`),
	}

	positions := mapPositions(columns, [][]json.RawMessage{row})
	if len(positions) != 1 {
		t.Fatalf("len = %d", len(positions))
	}
	p := positions[0]
	if p.ID != "360287970189927922" || p.TradableInstrumentID != 6109 || p.Side != "buy" {
		t.Fatalf("position = %+v", p)
	}
	if p.Fields["unrealizedPl"] != "-0.10" {
		t.Fatalf("unrealizedPl = %v", p.Fields["unrealizedPl"])
	}
}

func TestFilterPositionsByInstrument(t *testing.T) {
	positions := []PositionSummary{
		{ID: "1", TradableInstrumentID: 100},
		{ID: "2", TradableInstrumentID: 200},
	}
	filtered := FilterPositionsByInstrument(positions, 100)
	if len(filtered) != 1 || filtered[0].ID != "1" {
		t.Fatalf("filtered = %+v", filtered)
	}
}

func TestMapInstruments(t *testing.T) {
	instruments := []Instrument{
		{
			ID:                   1,
			Name:                 "BTCUSD",
			TradableInstrumentID: 2001,
			Type:                 "CRYPTO",
			TradingExchange:      "CRYPTO",
			MarketDataExchange:   "CRYPTO",
			Description:          "Bitcoin vs USD",
			Routes:               []Route{{ID: 5, Type: "TRADE"}, {ID: 6, Type: "INFO"}},
		},
	}

	summaries := mapInstruments(instruments)
	if len(summaries) != 1 {
		t.Fatalf("len = %d", len(summaries))
	}
	s := summaries[0]
	if s.Name != "BTCUSD" || s.TradeRouteID != 5 || s.InfoRouteID != 6 {
		t.Fatalf("summary = %+v", s)
	}
}
