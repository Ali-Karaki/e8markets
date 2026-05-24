package tradelocker

import (
	"context"
	"net/http"
	"time"
)

func (c *Client) Login(ctx context.Context, sessionID *string, email, password, server string) (Tokens, error) {
	req := LoginRequest{
		Email:    email,
		Password: password,
		Server:   server,
	}

	var resp TokenResponse
	_, err := c.doJSON(ctx, sessionID, http.MethodPost, pathAuthToken, "", req, "", &resp)
	if err != nil {
		return Tokens{}, err
	}

	expiresAt, err := time.Parse(time.RFC3339, resp.ExpireDate)
	if err != nil {
		expiresAt, err = time.Parse(time.RFC3339Nano, resp.ExpireDate)
		if err != nil {
			return Tokens{}, err
		}
	}

	return Tokens{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (c *Client) Refresh(ctx context.Context, sessionID *string, accessToken, refreshToken string) (Tokens, error) {
	req := RefreshRequest{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	var resp TokenResponse
	_, err := c.doJSON(ctx, sessionID, http.MethodPost, pathAuthRefresh, "", req, "", &resp)
	if err != nil {
		return Tokens{}, err
	}

	expiresAt, err := time.Parse(time.RFC3339, resp.ExpireDate)
	if err != nil {
		expiresAt, err = time.Parse(time.RFC3339Nano, resp.ExpireDate)
		if err != nil {
			return Tokens{}, err
		}
	}

	return Tokens{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (c *Client) GetAllAccounts(ctx context.Context, sessionID *string, accessToken string) ([]Account, error) {
	var resp AccountsResponse
	_, err := c.doJSON(ctx, sessionID, http.MethodGet, pathAuthAllAccounts, "", nil, accessToken, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Accounts, nil
}
