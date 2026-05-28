package db

import (
	"context"
	"fmt"

	"github.com/ArnobKumarSaha/football_calc/pkg/models"
	"github.com/ArnobKumarSaha/football_calc/pkg/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPlayerBalance(ctx context.Context, pool *pgxpool.Pool, playerID int) (*models.BalanceHistory, error) {
	player, err := GetPlayer(ctx, pool, playerID)
	if err != nil {
		return nil, fmt.Errorf("get player: %w", err)
	}

	rows, err := pool.Query(ctx,
		`SELECT ma.match_id, m.date::text,
		   COALESCE(SUM(pay.amount), 0) AS paid,
		   (SELECT MAX(id) FROM match_attendees WHERE match_id=m.id) AS last_id,
		   (SELECT SUM(1+guest_count) FROM match_attendees WHERE match_id=m.id) AS total_units,
		   ma.guest_count, ma.id AS my_id, m.total_bill
		 FROM match_attendees ma
		 JOIN matches m ON m.id=ma.match_id AND m.deleted_at IS NULL
		 LEFT JOIN payments pay ON pay.player_id=ma.player_id AND pay.match_id=ma.match_id
		 WHERE ma.player_id=$1
		 GROUP BY ma.match_id, m.id, m.date, ma.guest_count, ma.id, m.total_bill
		 ORDER BY m.date ASC`, playerID)
	if err != nil {
		return nil, fmt.Errorf("balance history: %w", err)
	}
	defer rows.Close()

	h := &models.BalanceHistory{
		PlayerID:   player.ID,
		PlayerName: player.Name,
		Matches:    []models.MatchBalanceEntry{},
	}
	for rows.Next() {
		var entry models.MatchBalanceEntry
		var lastID, myID, guestCount int
		var totalUnits, totalBill float64
		if err := rows.Scan(&entry.MatchID, &entry.MatchDate, &entry.Paid, &lastID, &totalUnits, &guestCount, &myID, &totalBill); err != nil {
			return nil, err
		}
		if myID == lastID {
			others, err := sumOtherDues(ctx, pool, entry.MatchID, lastID, totalBill, totalUnits)
			if err != nil {
				return nil, err
			}
			entry.Due = totalBill - others
		} else {
			entry.Due = service.Round2(totalBill / totalUnits * float64(1+guestCount))
		}
		h.Matches = append(h.Matches, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	allPayments, err := ListPayments(ctx, pool, &playerID, nil)
	if err != nil {
		return nil, err
	}
	h.Standalone = []models.Payment{}
	for _, p := range allPayments {
		if p.MatchID == nil {
			h.Standalone = append(h.Standalone, p)
		}
	}

	totalDue, err := playerTotalDue(ctx, pool, playerID)
	if err != nil {
		return nil, err
	}
	totalPaid, err := playerTotalPaid(ctx, pool, playerID)
	if err != nil {
		return nil, err
	}
	h.TotalDue = totalDue
	h.TotalPaid = totalPaid
	h.Balance = totalPaid - totalDue
	return h, nil
}

func AllPlayerBalances(ctx context.Context, pool *pgxpool.Pool) ([]models.PlayerWithBalance, error) {
	rows, err := pool.Query(ctx,
		`SELECT id, name, created_at FROM players WHERE deleted_at IS NULL ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("all balances: %w", err)
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var p models.Player
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]models.PlayerWithBalance, 0, len(players))
	for _, p := range players {
		due, err := playerTotalDue(ctx, pool, p.ID)
		if err != nil {
			return nil, err
		}
		paid, err := playerTotalPaid(ctx, pool, p.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, models.PlayerWithBalance{
			Player:  p,
			Balance: paid - due,
		})
	}
	return result, nil
}

func playerTotalPaid(ctx context.Context, pool *pgxpool.Pool, playerID int) (float64, error) {
	var total float64
	err := pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM payments WHERE player_id=$1`, playerID).
		Scan(&total)
	return total, err
}

func playerTotalDue(ctx context.Context, pool *pgxpool.Pool, playerID int) (float64, error) {
	rows, err := pool.Query(ctx,
		`SELECT m.id, m.total_bill,
		   (SELECT MAX(id) FROM match_attendees WHERE match_id=m.id) AS last_id,
		   (SELECT SUM(1+guest_count) FROM match_attendees WHERE match_id=m.id) AS total_units,
		   ma.guest_count, ma.id
		 FROM match_attendees ma
		 JOIN matches m ON m.id=ma.match_id AND m.deleted_at IS NULL
		 WHERE ma.player_id=$1`, playerID)
	if err != nil {
		return 0, fmt.Errorf("player due: %w", err)
	}
	defer rows.Close()

	var totalDue float64
	for rows.Next() {
		var matchID, lastID, guestCount, myID int
		var totalBill, totalUnits float64
		if err := rows.Scan(&matchID, &totalBill, &lastID, &totalUnits, &guestCount, &myID); err != nil {
			return 0, err
		}
		if myID == lastID {
			others, err := sumOtherDues(ctx, pool, matchID, lastID, totalBill, totalUnits)
			if err != nil {
				return 0, err
			}
			totalDue += totalBill - others
		} else {
			perUnit := totalBill / totalUnits
			totalDue += service.Round2(perUnit * float64(1+guestCount))
		}
	}
	return totalDue, rows.Err()
}

func sumOtherDues(ctx context.Context, pool *pgxpool.Pool, matchID, lastID int, totalBill, totalUnits float64) (float64, error) {
	rows, err := pool.Query(ctx,
		`SELECT guest_count FROM match_attendees WHERE match_id=$1 AND id!=$2`, matchID, lastID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	perUnit := totalBill / totalUnits
	var sum float64
	for rows.Next() {
		var gc int
		if err := rows.Scan(&gc); err != nil {
			return 0, err
		}
		sum += service.Round2(perUnit * float64(1+gc))
	}
	return sum, rows.Err()
}
