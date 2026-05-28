CREATE TABLE IF NOT EXISTS players (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS matches (
    id         SERIAL PRIMARY KEY,
    date       DATE NOT NULL,
    total_bill NUMERIC(10,2) NOT NULL,
    notes      TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS match_attendees (
    id          SERIAL PRIMARY KEY,
    match_id    INT NOT NULL REFERENCES matches(id),
    player_id   INT NOT NULL REFERENCES players(id),
    guest_count INT NOT NULL DEFAULT 0,
    UNIQUE (match_id, player_id)
);

CREATE TABLE IF NOT EXISTS payments (
    id         SERIAL PRIMARY KEY,
    player_id  INT NOT NULL REFERENCES players(id),
    match_id   INT REFERENCES matches(id),
    amount     NUMERIC(10,2) NOT NULL,
    notes      TEXT,
    paid_at    DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Drop old full-table unique on players.name if present; replaced by partial index below.
ALTER TABLE players DROP CONSTRAINT IF EXISTS players_name_key;

-- Only active (non-deleted) players must have unique names.
CREATE UNIQUE INDEX IF NOT EXISTS players_name_active ON players(name) WHERE deleted_at IS NULL;

-- Indexes on FK columns not covered by an existing index.
CREATE INDEX IF NOT EXISTS payments_player_id_idx ON payments(player_id);
CREATE INDEX IF NOT EXISTS payments_match_id_idx  ON payments(match_id);
CREATE INDEX IF NOT EXISTS attendees_player_id_idx ON match_attendees(player_id);

-- CHECK constraints (idempotent: no-op if already present).
DO $$ BEGIN
    ALTER TABLE matches ADD CONSTRAINT matches_total_bill_positive CHECK (total_bill > 0);
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE payments DROP CONSTRAINT IF EXISTS payments_amount_positive;
    ALTER TABLE payments ADD CONSTRAINT payments_amount_nonzero CHECK (amount <> 0);
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
    ALTER TABLE match_attendees ADD CONSTRAINT attendees_guest_count_nonneg CHECK (guest_count >= 0);
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;
