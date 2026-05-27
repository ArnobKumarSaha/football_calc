package db

import (
	"context"
	"fmt"

	"github.com/ArnobKumarSaha/football_calc/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ListMatches(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]models.Match, error) {
	rows, err := pool.Query(ctx,
		`SELECT id, date::text, total_bill, notes, created_at, deleted_at
		 FROM matches WHERE deleted_at IS NULL ORDER BY date DESC LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list matches: %w", err)
	}
	defer rows.Close()

	var matches []models.Match
	for rows.Next() {
		var m models.Match
		if err := rows.Scan(&m.ID, &m.Date, &m.TotalBill, &m.Notes, &m.CreatedAt, &m.DeletedAt); err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

func GetMatch(ctx context.Context, pool *pgxpool.Pool, id int) (*models.Match, error) {
	var m models.Match
	err := pool.QueryRow(ctx,
		`SELECT id, date::text, total_bill, notes, created_at, deleted_at
		 FROM matches WHERE id=$1 AND deleted_at IS NULL`, id).
		Scan(&m.ID, &m.Date, &m.TotalBill, &m.Notes, &m.CreatedAt, &m.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func CreateMatch(ctx context.Context, pool *pgxpool.Pool, date string, totalBill float64, notes *string) (*models.Match, error) {
	var m models.Match
	err := pool.QueryRow(ctx,
		`INSERT INTO matches (date, total_bill, notes) VALUES ($1,$2,$3)
		 RETURNING id, date::text, total_bill, notes, created_at, deleted_at`,
		date, totalBill, notes).
		Scan(&m.ID, &m.Date, &m.TotalBill, &m.Notes, &m.CreatedAt, &m.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("create match: %w", err)
	}
	return &m, nil
}

func UpdateMatch(ctx context.Context, pool *pgxpool.Pool, id int, date *string, totalBill *float64, notes *string) (*models.Match, error) {
	var m models.Match
	err := pool.QueryRow(ctx,
		`UPDATE matches SET
		   date       = COALESCE($2, date),
		   total_bill = COALESCE($3, total_bill),
		   notes      = COALESCE($4, notes)
		 WHERE id=$1 AND deleted_at IS NULL
		 RETURNING id, date::text, total_bill, notes, created_at, deleted_at`,
		id, date, totalBill, notes).
		Scan(&m.ID, &m.Date, &m.TotalBill, &m.Notes, &m.CreatedAt, &m.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("update match: %w", err)
	}
	return &m, nil
}

func DeleteMatch(ctx context.Context, pool *pgxpool.Pool, id int) error {
	cmd, err := pool.Exec(ctx,
		`UPDATE matches SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return fmt.Errorf("delete match: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("match not found")
	}
	return nil
}
