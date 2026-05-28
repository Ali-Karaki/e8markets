package handlers

import (
	"log"
	"net/http"

	"github.com/Ali-Karaki/e8markets/server/internal/httpx"
	"github.com/Ali-Karaki/e8markets/server/internal/middleware"
	"github.com/Ali-Karaki/e8markets/server/internal/tradelocker"
)

type InstrumentsHandler struct {
	tl *tradelocker.Client
}

type instrumentsListResponse struct {
	Instruments []tradelocker.InstrumentSummary `json:"instruments"`
}

func NewInstrumentsHandler(tl *tradelocker.Client) *InstrumentsHandler {
	return &InstrumentsHandler{tl: tl}
}

func (h *InstrumentsHandler) List(w http.ResponseWriter, r *http.Request) {
	session, ok := middleware.SessionFromContext(r.Context())
	if !ok {
		return
	}

	accountID := r.URL.Query().Get("accountId")
	accNum := r.URL.Query().Get("accNum")
	if accountID == "" || accNum == "" {
		writeValidationError(w, "accountId and accNum are required")
		return
	}

	sid := session.ID.String()
	instruments, err := h.tl.GetInstruments(r.Context(), &sid, session.AccessToken, accountID, accNum)
	if err != nil {
		log.Printf("GetInstruments failed session=%s accountId=%s accNum=%s err=%v", sid, accountID, accNum, err)
		writeTradeLockerError(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, instrumentsListResponse{Instruments: instruments})
}
