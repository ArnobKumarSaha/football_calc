CREATE TABLE IF NOT EXISTS players (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
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
