package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func HasAnyData(ctx context.Context, pool *pgxpool.Pool) (bool, error) {
	var count int
	err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM players`).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check data: %w", err)
	}
	return count > 0, nil
}

type ImportEntry struct {
	Name    string
	Balance float64
}

func ImportInitialData(ctx context.Context, pool *pgxpool.Pool, entries []ImportEntry) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, e := range entries {
		var playerID int
		err := tx.QueryRow(ctx,
			`INSERT INTO players (name) VALUES ($1) RETURNING id`, e.Name).Scan(&playerID)
		if err != nil {
			return fmt.Errorf("insert player %q: %w", e.Name, err)
		}

		if e.Balance != 0 {
			_, err = tx.Exec(ctx,
				`INSERT INTO payments (player_id, amount, paid_at, notes) VALUES ($1, $2, CURRENT_DATE, 'initial balance')`,
				playerID, e.Balance)
			if err != nil {
				return fmt.Errorf("insert initial payment for %q: %w", e.Name, err)
			}
		}
	}

	return tx.Commit(ctx)
}
