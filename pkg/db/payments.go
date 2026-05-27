package db

import (
	"context"
	"fmt"

	"github.com/ArnobKumarSaha/football_calc/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ListPayments(ctx context.Context, pool *pgxpool.Pool, playerID, matchID *int) ([]models.Payment, error) {
	query := `SELECT id, player_id, match_id, amount, notes, paid_at, created_at FROM payments WHERE 1=1`
	args := []any{}
	idx := 1

	if playerID != nil {
		query += fmt.Sprintf(" AND player_id=$%d", idx)
		args = append(args, *playerID)
		idx++
	}
	if matchID != nil {
		query += fmt.Sprintf(" AND match_id=$%d", idx)
		args = append(args, *matchID)
		idx++
	}
	query += " ORDER BY paid_at DESC, id DESC"

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list payments: %w", err)
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(&p.ID, &p.PlayerID, &p.MatchID, &p.Amount, &p.Notes, &p.PaidAt, &p.CreatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, rows.Err()
}

func CreatePayment(ctx context.Context, pool *pgxpool.Pool, playerID int, matchID *int, amount float64, paidAt string, notes *string) (*models.Payment, error) {
	var p models.Payment
	err := pool.QueryRow(ctx,
		`INSERT INTO payments (player_id, match_id, amount, paid_at, notes) VALUES ($1,$2,$3,$4,$5)
		 RETURNING id, player_id, match_id, amount, notes, paid_at, created_at`,
		playerID, matchID, amount, paidAt, notes).
		Scan(&p.ID, &p.PlayerID, &p.MatchID, &p.Amount, &p.Notes, &p.PaidAt, &p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}
	return &p, nil
}

func UpdatePayment(ctx context.Context, pool *pgxpool.Pool, id int, matchID *int, amount *float64, paidAt *string, notes *string) (*models.Payment, error) {
	var p models.Payment
	err := pool.QueryRow(ctx,
		`UPDATE payments SET
		   match_id = COALESCE($2, match_id),
		   amount   = COALESCE($3, amount),
		   paid_at  = COALESCE($4, paid_at),
		   notes    = COALESCE($5, notes)
		 WHERE id=$1
		 RETURNING id, player_id, match_id, amount, notes, paid_at, created_at`,
		id, matchID, amount, paidAt, notes).
		Scan(&p.ID, &p.PlayerID, &p.MatchID, &p.Amount, &p.Notes, &p.PaidAt, &p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("update payment: %w", err)
	}
	return &p, nil
}

func DeletePayment(ctx context.Context, pool *pgxpool.Pool, id int) error {
	cmd, err := pool.Exec(ctx, `DELETE FROM payments WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete payment: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("payment not found")
	}
	return nil
}
