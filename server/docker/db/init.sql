CREATE TABLE IF NOT EXISTS api_logs (
  id            BIGSERIAL PRIMARY KEY,
  session_id    UUID,
  event_type    TEXT NOT NULL,
  method        TEXT,
  path          TEXT,
  status_code   INT,
  message       TEXT,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sync_runs (
  id              BIGSERIAL PRIMARY KEY,
  session_id      UUID NOT NULL,
  account_id      TEXT NOT NULL,
  acc_num         TEXT NOT NULL,
  sync_type       TEXT NOT NULL DEFAULT 'positions',
  status          TEXT NOT NULL,
  records_stored  INT NOT NULL DEFAULT 0,
  error_message   TEXT,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_sync_runs_account
  ON sync_runs (account_id, acc_num, created_at DESC);

CREATE TABLE IF NOT EXISTS position_snapshots (
  id                      BIGSERIAL PRIMARY KEY,
  sync_run_id             BIGINT NOT NULL REFERENCES sync_runs(id) ON DELETE CASCADE,
  session_id              UUID NOT NULL,
  account_id              TEXT NOT NULL,
  acc_num                 TEXT NOT NULL,
  position_id             TEXT NOT NULL,
  tradable_instrument_id  BIGINT,
  data                    JSONB NOT NULL,
  created_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_position_snapshots_history
  ON position_snapshots (account_id, acc_num, tradable_instrument_id, created_at DESC);
