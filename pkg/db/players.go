package db

import (
	"context"
	"fmt"

	"github.com/ArnobKumarSaha/football_calc/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPlayer(ctx context.Context, pool *pgxpool.Pool, id int) (*models.Player, error) {
	var p models.Player
	err := pool.QueryRow(ctx,
		`SELECT id, name, created_at, deleted_at FROM players WHERE id=$1 AND deleted_at IS NULL`, id).
		Scan(&p.ID, &p.Name, &p.CreatedAt, &p.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func CreatePlayer(ctx context.Context, pool *pgxpool.Pool, name string) (*models.Player, error) {
	var p models.Player
	err := pool.QueryRow(ctx,
		`INSERT INTO players (name) VALUES ($1) RETURNING id, name, created_at, deleted_at`, name).
		Scan(&p.ID, &p.Name, &p.CreatedAt, &p.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("create player: %w", err)
	}
	return &p, nil
}

func DeletePlayer(ctx context.Context, pool *pgxpool.Pool, id int) error {
	cmd, err := pool.Exec(ctx,
		`UPDATE players SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return fmt.Errorf("delete player: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("player not found")
	}
	return nil
}
