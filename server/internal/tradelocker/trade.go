package tradelocker

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) GetAccountState(ctx context.Context, sessionID *string, accessToken, accountID, accNum string) (map[string]float64, error) {
	if err := c.loadConfig(ctx, sessionID, accessToken, accNum); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(pathTradeAccountStateFmt, accountID)
	var resp TLResponse[AccountStateData]
	_, err := c.doJSON(ctx, sessionID, http.MethodGet, path, accNum, nil, accessToken, &resp)
	if err != nil {
		return nil, err
	}

	return mapAccountDetails(c.columns, resp.D.AccountDetailsData), nil
}

func mapAccountDetails(columns []Column, data []flexFloat64) map[string]float64 {
	result := make(map[string]float64)
	for i, val := range data {
		if i < len(columns) {
			key := columns[i].ID
			if key != "" {
				result[key] = val.Float64()
			}
		}
	}
	return result
}

func (c *Client) loadConfig(ctx context.Context, sessionID *string, accessToken, accNum string) error {
	c.configOnce.Do(func() {
		var resp TLResponse[ConfigData]
		_, err := c.doJSON(ctx, sessionID, http.MethodGet, pathTradeConfig, accNum, nil, accessToken, &resp)
		if err != nil {
			c.configErr = err
			return
		}
		if resp.D.AccountDetailsConfig != nil {
			c.columns = resp.D.AccountDetailsConfig.Columns
		}
	})
	return c.configErr
}
