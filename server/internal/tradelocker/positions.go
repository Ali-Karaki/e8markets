package tradelocker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (c *Client) GetPositions(ctx context.Context, sessionID *string, accessToken, accountID, accNum string) ([]PositionSummary, error) {
	if err := c.loadConfig(ctx, sessionID, accessToken, accNum); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(pathTradeAccountPositionsFmt, accountID)
	var resp TLResponse[PositionsData]
	_, err := c.doJSON(ctx, sessionID, http.MethodGet, path, accNum, nil, accessToken, &resp)
	if err != nil {
		return nil, err
	}

	return mapPositions(c.positionColumns, resp.D.Positions), nil
}

func mapPositions(columns []Column, rows [][]json.RawMessage) []PositionSummary {
	result := make([]PositionSummary, 0, len(rows))
	for _, row := range rows {
		result = append(result, mapPositionRow(columns, row))
	}
	return result
}

func mapPositionRow(columns []Column, row []json.RawMessage) PositionSummary {
	fields := make(map[string]string)
	for i, cell := range row {
		if i >= len(columns) {
			break
		}
		key := columns[i].ID
		if key == "" {
			continue
		}
		fields[key] = parseCell(cell)
	}

	summary := PositionSummary{Fields: fields}
	if v, ok := fields["id"]; ok {
		summary.ID = v
	}
	if v, ok := fields["tradableInstrumentId"]; ok {
		summary.TradableInstrumentID = parseInt64(v)
	}
	if v, ok := fields["side"]; ok {
		summary.Side = v
	}
	if v, ok := fields["qty"]; ok {
		summary.Qty = parseFloat64(v)
	}
	if v, ok := fields["avgPrice"]; ok {
		summary.AvgPrice = parseFloat64(v)
	}
	if v, ok := fields["routeId"]; ok {
		summary.RouteID = parseInt64(v)
	}
	if v, ok := fields["stopLossId"]; ok {
		summary.StopLossID = parseInt64(v)
	}
	if v, ok := fields["takeProfitId"]; ok {
		summary.TakeProfitID = parseInt64(v)
	}

	if len(summary.Fields) == 0 {
		summary.Fields = nil
	}

	return summary
}

func parseCell(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}

	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}

	var n float64
	if err := json.Unmarshal(raw, &n); err == nil {
		return strconv.FormatFloat(n, 'f', -1, 64)
	}

	return string(raw)
}

func parseFloat64(s string) float64 {
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func parseInt64(s string) int64 {
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

// HasNewOpenPositions reports whether any open position ID is not in previousSyncedIDs.
// An empty previous map means never synced: any current open position counts as new.
func HasNewOpenPositions(current []PositionSummary, previousSyncedIDs map[string]struct{}) bool {
	for _, p := range current {
		if p.ID == "" {
			continue
		}
		if _, seen := previousSyncedIDs[p.ID]; !seen {
			return true
		}
	}
	return false
}

func FilterPositionsByInstrument(positions []PositionSummary, tradableInstrumentID int64) []PositionSummary {
	if tradableInstrumentID == 0 {
		return positions
	}
	filtered := make([]PositionSummary, 0, len(positions))
	for _, p := range positions {
		if p.TradableInstrumentID == tradableInstrumentID {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
