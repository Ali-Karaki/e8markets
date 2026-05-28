package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Ali-Karaki/e8markets/server/internal/httpx"
	"github.com/Ali-Karaki/e8markets/server/internal/middleware"
	"github.com/Ali-Karaki/e8markets/server/internal/store"
	"github.com/Ali-Karaki/e8markets/server/internal/tradelocker"
)

type PositionsHandler struct {
	tl    *tradelocker.Client
	store *store.PositionStore
}

func NewPositionsHandler(tl *tradelocker.Client, positionStore *store.PositionStore) *PositionsHandler {
	return &PositionsHandler{tl: tl, store: positionStore}
}

type positionsListResponse struct {
	Positions []tradelocker.PositionSummary `json:"positions"`
}

type syncPositionsResponse struct {
	SyncRunID     int64                           `json:"syncRunId"`
	Status        string                          `json:"status"`
	RecordsStored int                             `json:"recordsStored"`
	Skipped       bool                            `json:"skipped,omitempty"`
	Positions     []tradelocker.PositionSummary   `json:"positions"`
}

type positionHistoryItem struct {
	ID         int64                         `json:"id"`
	SyncRunID  int64                         `json:"syncRunId"`
	SyncedAt   string                        `json:"syncedAt"`
	PositionID string                        `json:"positionId"`
	Position   tradelocker.PositionSummary   `json:"position"`
}

type positionHistoryResponse struct {
	Snapshots []positionHistoryItem `json:"snapshots"`
}

func (h *PositionsHandler) List(w http.ResponseWriter, r *http.Request) {
	session, ok := middleware.SessionFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	accountID, accNum, instrumentID, ok := parseAccountQuery(w, r)
	if !ok {
		return
	}

	sid := session.ID.String()
	positions, err := h.tl.GetPositions(r.Context(), &sid, session.AccessToken, accountID, accNum)
	if err != nil {
		log.Printf("GetPositions failed session=%s accountId=%s accNum=%s err=%v", sid, accountID, accNum, err)
		httpx.Error(w, http.StatusBadGateway, "Failed to fetch positions")
		return
	}

	if instrumentID != nil {
		positions = tradelocker.FilterPositionsByInstrument(positions, *instrumentID)
	}

	httpx.JSON(w, http.StatusOK, positionsListResponse{Positions: positions})
}

func (h *PositionsHandler) Sync(w http.ResponseWriter, r *http.Request) {
	session, ok := middleware.SessionFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	accountID, accNum, _, ok := parseAccountQuery(w, r)
	if !ok {
		return
	}

	sid := session.ID.String()
	positions, err := h.tl.GetPositions(r.Context(), &sid, session.AccessToken, accountID, accNum)
	if err != nil {
		log.Printf("Sync GetPositions failed session=%s accountId=%s accNum=%s err=%v", sid, accountID, accNum, err)
		_, _ = h.store.RecordFailedSync(r.Context(), session.ID, accountID, accNum, "Failed to fetch positions from TradeLocker")
		httpx.Error(w, http.StatusBadGateway, "Failed to fetch positions")
		return
	}

	force, onlyIfNew := parseSyncOptions(r)
	if onlyIfNew && !force {
		previous, err := h.store.GetOpenPositionIDsFromLatestSync(r.Context(), accountID, accNum)
		if err != nil {
			log.Printf("GetOpenPositionIDsFromLatestSync failed accountId=%s accNum=%s err=%v", accountID, accNum, err)
			httpx.Error(w, http.StatusInternalServerError, "Failed to check position sync state")
			return
		}
		if !tradelocker.HasNewOpenPositions(positions, previous) {
			httpx.JSON(w, http.StatusOK, syncPositionsResponse{
				Status:        store.SyncStatusSkipped,
				RecordsStored: 0,
				Skipped:       true,
				Positions:     positions,
			})
			return
		}
	}

	inputs := make([]store.SnapshotInput, 0, len(positions))
	for _, pos := range positions {
		data, err := json.Marshal(pos)
		if err != nil {
			log.Printf("marshal position session=%s err=%v", sid, err)
			httpx.Error(w, http.StatusInternalServerError, "Failed to store positions")
			return
		}

		var tradableID *int64
		if pos.TradableInstrumentID != 0 {
			id := pos.TradableInstrumentID
			tradableID = &id
		}

		inputs = append(inputs, store.SnapshotInput{
			PositionID:           pos.ID,
			TradableInstrumentID: tradableID,
			Data:                 data,
		})
	}

	syncRunID, recordsStored, err := h.store.SyncPositions(r.Context(), session.ID, accountID, accNum, inputs)
	if err != nil {
		log.Printf("SyncPositions failed session=%s accountId=%s err=%v", sid, accountID, err)
		_, _ = h.store.RecordFailedSync(r.Context(), session.ID, accountID, accNum, "Failed to store position snapshots")
		httpx.Error(w, http.StatusInternalServerError, "Failed to store positions")
		return
	}

	httpx.JSON(w, http.StatusOK, syncPositionsResponse{
		SyncRunID:     syncRunID,
		Status:        store.SyncStatusSuccess,
		RecordsStored: recordsStored,
		Positions:     positions,
	})
}

func parseSyncOptions(r *http.Request) (force, onlyIfNew bool) {
	return r.URL.Query().Get("force") == "true", r.URL.Query().Get("onlyIfNew") == "true"
}

func (h *PositionsHandler) History(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.SessionFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	accountID, accNum, instrumentID, ok := parseAccountQuery(w, r)
	if !ok {
		return
	}

	snapshots, err := h.store.ListSnapshots(r.Context(), accountID, accNum, instrumentID, 0)
	if err != nil {
		log.Printf("ListSnapshots failed accountId=%s accNum=%s err=%v", accountID, accNum, err)
		httpx.Error(w, http.StatusInternalServerError, "Failed to fetch position history")
		return
	}

	items := make([]positionHistoryItem, 0, len(snapshots))
	for _, snap := range snapshots {
		var pos tradelocker.PositionSummary
		if err := json.Unmarshal(snap.Data, &pos); err != nil {
			log.Printf("unmarshal snapshot id=%d err=%v", snap.ID, err)
			continue
		}
		items = append(items, positionHistoryItem{
			ID:         snap.ID,
			SyncRunID:  snap.SyncRunID,
			SyncedAt:   snap.SyncedAt.UTC().Format(time.RFC3339),
			PositionID: snap.PositionID,
			Position:   pos,
		})
	}

	httpx.JSON(w, http.StatusOK, positionHistoryResponse{Snapshots: items})
}

func parseAccountQuery(w http.ResponseWriter, r *http.Request) (accountID, accNum string, instrumentID *int64, ok bool) {
	accountID = r.URL.Query().Get("accountId")
	accNum = r.URL.Query().Get("accNum")
	if accountID == "" || accNum == "" {
		httpx.Error(w, http.StatusBadRequest, "accountId and accNum are required")
		return "", "", nil, false
	}

	if raw := r.URL.Query().Get("tradableInstrumentId"); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "tradableInstrumentId must be a number")
			return "", "", nil, false
		}
		instrumentID = &id
	}

	return accountID, accNum, instrumentID, true
}
