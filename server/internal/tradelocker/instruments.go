package tradelocker

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) GetInstruments(ctx context.Context, sessionID *string, accessToken, accountID, accNum string) ([]InstrumentSummary, error) {
	path := fmt.Sprintf(pathTradeAccountInstrumentsFmt, accountID)
	var resp TLResponse[InstrumentsData]
	_, err := c.doJSON(ctx, sessionID, http.MethodGet, path, accNum, nil, accessToken, &resp)
	if err != nil {
		return nil, err
	}

	return mapInstruments(resp.D.Instruments), nil
}

func mapInstruments(instruments []Instrument) []InstrumentSummary {
	result := make([]InstrumentSummary, 0, len(instruments))
	for _, inst := range instruments {
		tradeRouteID, infoRouteID := extractRouteIDs(inst.Routes)
		result = append(result, InstrumentSummary{
			ID:                   inst.ID,
			TradableInstrumentID: inst.TradableInstrumentID,
			Name:                 inst.Name,
			Description:          inst.Description,
			Type:                 inst.Type,
			TradingExchange:      inst.TradingExchange,
			MarketDataExchange:   inst.MarketDataExchange,
			TradeRouteID:         tradeRouteID,
			InfoRouteID:          infoRouteID,
		})
	}
	return result
}

func extractRouteIDs(routes []Route) (tradeRouteID, infoRouteID int64) {
	for _, route := range routes {
		switch route.Type {
		case "TRADE":
			tradeRouteID = route.ID
		case "INFO":
			infoRouteID = route.ID
		}
	}
	return tradeRouteID, infoRouteID
}
