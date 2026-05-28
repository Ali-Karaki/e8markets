package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	SyncStatusSuccess = "success"
	SyncStatusFailure = "failure"
	SyncStatusSkipped = "skipped"
	SyncTypePositions = "positions"

	defaultHistoryLimit = 100
)

type PositionStore struct {
	pg *Postgres
}

func NewPositionStore(pg *Postgres) *PositionStore {
	return &PositionStore{pg: pg}
}

type SnapshotInput struct {
	PositionID           string
	TradableInstrumentID *int64
	Data                 json.RawMessage
}

type StoredSnapshot struct {
	ID         int64
	SyncRunID  int64
	SyncedAt   time.Time
	PositionID string
	Data       json.RawMessage
}

func (s *PositionStore) SyncPositions(
	ctx context.Context,
	sessionID uuid.UUID,
	accountID, accNum string,
	snapshots []SnapshotInput,
) (syncRunID int64, recordsStored int, err error) {
	tx, err := s.pg.pool.Begin(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO sync_runs (session_id, account_id, acc_num, sync_type, status, records_stored)
		VALUES ($1, $2, $3, $4, $5, 0)
		RETURNING id
	`, sessionID, accountID, accNum, SyncTypePositions, SyncStatusSuccess).Scan(&syncRunID)
	if err != nil {
		return 0, 0, fmt.Errorf("insert sync_run: %w", err)
	}

	for _, snap := range snapshots {
		_, err = tx.Exec(ctx, `
			INSERT INTO position_snapshots (
				sync_run_id, session_id, account_id, acc_num,
				position_id, tradable_instrument_id, data
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, syncRunID, sessionID, accountID, accNum, snap.PositionID, snap.TradableInstrumentID, snap.Data)
		if err != nil {
			return 0, 0, fmt.Errorf("insert snapshot: %w", err)
		}
		recordsStored++
	}

	_, err = tx.Exec(ctx, `
		UPDATE sync_runs SET records_stored = $1 WHERE id = $2
	`, recordsStored, syncRunID)
	if err != nil {
		return 0, 0, fmt.Errorf("update sync_run: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, 0, fmt.Errorf("commit tx: %w", err)
	}

	return syncRunID, recordsStored, nil
}

func (s *PositionStore) GetOpenPositionIDsFromLatestSync(
	ctx context.Context,
	accountID, accNum string,
) (map[string]struct{}, error) {
	rows, err := s.pg.pool.Query(ctx, `
		SELECT ps.position_id
		FROM position_snapshots ps
		WHERE ps.sync_run_id = (
			SELECT id FROM sync_runs
			WHERE account_id = $1 AND acc_num = $2 AND status = $3
			ORDER BY created_at DESC
			LIMIT 1
		)
	`, accountID, accNum, SyncStatusSuccess)
	if err != nil {
		return nil, fmt.Errorf("query latest sync position ids: %w", err)
	}
	defer rows.Close()

	ids := make(map[string]struct{})
	for rows.Next() {
		var positionID string
		if err := rows.Scan(&positionID); err != nil {
			return nil, fmt.Errorf("scan position_id: %w", err)
		}
		ids[positionID] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate position ids: %w", err)
	}
	return ids, nil
}

func (s *PositionStore) RecordFailedSync(
	ctx context.Context,
	sessionID uuid.UUID,
	accountID, accNum, errMsg string,
) (syncRunID int64, err error) {
	err = s.pg.pool.QueryRow(ctx, `
		INSERT INTO sync_runs (session_id, account_id, acc_num, sync_type, status, records_stored, error_message)
		VALUES ($1, $2, $3, $4, $5, 0, $6)
		RETURNING id
	`, sessionID, accountID, accNum, SyncTypePositions, SyncStatusFailure, errMsg).Scan(&syncRunID)
	if err != nil {
		return 0, fmt.Errorf("insert failed sync_run: %w", err)
	}
	return syncRunID, nil
}

func (s *PositionStore) RecordSkippedSync(
	ctx context.Context,
	sessionID uuid.UUID,
	accountID, accNum string,
) (syncRunID int64, err error) {
	err = s.pg.pool.QueryRow(ctx, `
		INSERT INTO sync_runs (session_id, account_id, acc_num, sync_type, status, records_stored)
		VALUES ($1, $2, $3, $4, $5, 0)
		RETURNING id
	`, sessionID, accountID, accNum, SyncTypePositions, SyncStatusSkipped).Scan(&syncRunID)
	if err != nil {
		return 0, fmt.Errorf("insert skipped sync_run: %w", err)
	}
	return syncRunID, nil
}

func (s *PositionStore) ListSnapshots(
	ctx context.Context,
	accountID, accNum string,
	tradableInstrumentID *int64,
	limit int,
) ([]StoredSnapshot, error) {
	if limit <= 0 {
		limit = defaultHistoryLimit
	}

	var rows pgx.Rows
	var err error

	if tradableInstrumentID != nil {
		rows, err = s.pg.pool.Query(ctx, `
			SELECT id, sync_run_id, created_at, position_id, data
			FROM position_snapshots
			WHERE account_id = $1 AND acc_num = $2 AND tradable_instrument_id = $3
			ORDER BY created_at DESC
			LIMIT $4
		`, accountID, accNum, *tradableInstrumentID, limit)
	} else {
		rows, err = s.pg.pool.Query(ctx, `
			SELECT id, sync_run_id, created_at, position_id, data
			FROM position_snapshots
			WHERE account_id = $1 AND acc_num = $2
			ORDER BY created_at DESC
			LIMIT $3
		`, accountID, accNum, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("query snapshots: %w", err)
	}
	defer rows.Close()

	var result []StoredSnapshot
	for rows.Next() {
		var snap StoredSnapshot
		if err := rows.Scan(&snap.ID, &snap.SyncRunID, &snap.SyncedAt, &snap.PositionID, &snap.Data); err != nil {
			return nil, fmt.Errorf("scan snapshot: %w", err)
		}
		result = append(result, snap)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate snapshots: %w", err)
	}

	return result, nil
}
