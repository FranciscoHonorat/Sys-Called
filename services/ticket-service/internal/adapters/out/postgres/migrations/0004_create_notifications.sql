CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY,
    audience_kind TEXT NOT NULL,
    audience TEXT NOT NULL,
    message TEXT NOT NULL,
    ticket_id UUID,
    actor_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_notifications_audience ON notifications (audience_kind, audience, created_at DESC);

CREATE TABLE IF NOT EXISTS notification_reads (
    user_id TEXT PRIMARY KEY,
    seen_until TIMESTAMPTZ NOT NULL
);
