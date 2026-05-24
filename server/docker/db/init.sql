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
