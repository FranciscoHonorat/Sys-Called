CREATE TABLE IF NOT EXISTS responsibles (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    registered_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
