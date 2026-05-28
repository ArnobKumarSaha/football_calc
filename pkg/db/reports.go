package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MonthlyReport struct {
	Month      string  `json:"month"`
	MatchCount int     `json:"match_count"`
	TotalSpent float64 `json:"total_spent"`
}

type MatchCSVRow struct {
	MatchID    int
	Date       string
	TotalBill  float64
	PlayerName string
	GuestCount int
}

type PaymentCSVRow struct {
	ID         int
	PlayerID   int
	PlayerName string
	MatchID    *int
	Amount     float64
	PaidAt     string
	Notes      *string
}

func MonthlyReports(ctx context.Context, pool *pgxpool.Pool, year int) ([]MonthlyReport, error) {
	rows, err := pool.Query(ctx,
		`SELECT TO_CHAR(date,'YYYY-MM') AS month, COUNT(*) AS match_count, SUM(total_bill) AS total_spent
		 FROM matches WHERE deleted_at IS NULL AND EXTRACT(YEAR FROM date)=$1
		 GROUP BY month ORDER BY month`, year)
	if err != nil {
		return nil, fmt.Errorf("monthly reports: %w", err)
	}
	defer rows.Close()

	var reports []MonthlyReport
	for rows.Next() {
		var r MonthlyReport
		if err := rows.Scan(&r.Month, &r.MatchCount, &r.TotalSpent); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, rows.Err()
}

func MatchesCSV(ctx context.Context, pool *pgxpool.Pool) ([]MatchCSVRow, error) {
	rows, err := pool.Query(ctx,
		`SELECT m.id, m.date::text, m.total_bill, p.name, ma.guest_count
		 FROM matches m
		 JOIN match_attendees ma ON ma.match_id=m.id
		 JOIN players p ON p.id=ma.player_id
		 WHERE m.deleted_at IS NULL
		 ORDER BY m.date DESC, ma.id`)
	if err != nil {
		return nil, fmt.Errorf("matches csv: %w", err)
	}
	defer rows.Close()

	var result []MatchCSVRow
	for rows.Next() {
		var r MatchCSVRow
		if err := rows.Scan(&r.MatchID, &r.Date, &r.TotalBill, &r.PlayerName, &r.GuestCount); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, rows.Err()
}

func PaymentsCSV(ctx context.Context, pool *pgxpool.Pool) ([]PaymentCSVRow, error) {
	rows, err := pool.Query(ctx,
		`SELECT pay.id, p.id, p.name, pay.match_id, pay.amount, pay.paid_at::text, pay.notes
		 FROM payments pay
		 JOIN players p ON p.id=pay.player_id
		 ORDER BY pay.paid_at DESC, pay.id DESC`)
	if err != nil {
		return nil, fmt.Errorf("payments csv: %w", err)
	}
	defer rows.Close()

	var result []PaymentCSVRow
	for rows.Next() {
		var r PaymentCSVRow
		if err := rows.Scan(&r.ID, &r.PlayerID, &r.PlayerName, &r.MatchID, &r.Amount, &r.PaidAt, &r.Notes); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, rows.Err()
}
