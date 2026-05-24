package handlers

import (
	"net/http"

	"github.com/Ali-Karaki/e8markets/server/internal/httpx"
	"github.com/Ali-Karaki/e8markets/server/internal/middleware"
	"github.com/Ali-Karaki/e8markets/server/internal/tradelocker"
)

type AccountsHandler struct {
	tl *tradelocker.Client
}

func NewAccountsHandler(tl *tradelocker.Client) *AccountsHandler {
	return &AccountsHandler{tl: tl}
}

func (h *AccountsHandler) List(w http.ResponseWriter, r *http.Request) {
	session, ok := middleware.SessionFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	sid := session.ID.String()
	accounts, err := h.tl.GetAllAccounts(r.Context(), &sid, session.AccessToken)
	if err != nil {
		httpx.Error(w, http.StatusBadGateway, "Failed to fetch accounts")
		return
	}

	for i := range accounts {
		if accounts[i].AccountBalance == 0 && accounts[i].AAccountBalance != 0 {
			accounts[i].AccountBalance = accounts[i].AAccountBalance
		}
	}

	httpx.JSON(w, http.StatusOK, map[string]any{"accounts": accounts})
}

func (h *AccountsHandler) State(w http.ResponseWriter, r *http.Request) {
	session, ok := middleware.SessionFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	accountID := r.URL.Query().Get("accountId")
	accNum := r.URL.Query().Get("accNum")
	if accountID == "" || accNum == "" {
		httpx.Error(w, http.StatusBadRequest, "accountId and accNum are required")
		return
	}

	sid := session.ID.String()
	state, err := h.tl.GetAccountState(r.Context(), &sid, session.AccessToken, accountID, accNum)
	if err != nil {
		httpx.Error(w, http.StatusBadGateway, "Failed to fetch account state")
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]any{
		"accountId": accountID,
		"accNum":    accNum,
		"state":     state,
	})
}
