CREATE TABLE ticket_events (
    id UUID PRIMARY KEY,
    aggregate_id UUID NOT NULL,
    version INTEGER NOT NULL,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    UNIQUE (aggregate_id, version)
);

CREATE INDEX idx_ticket_events_aggregate_id ON ticket_events (aggregate_id, version);
