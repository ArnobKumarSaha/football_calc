package db

import (
	"context"
	"fmt"

	"github.com/ArnobKumarSaha/football_calc/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetMatchAttendees(ctx context.Context, pool *pgxpool.Pool, matchID int) ([]models.MatchAttendee, error) {
	rows, err := pool.Query(ctx,
		`SELECT id, match_id, player_id, guest_count FROM match_attendees WHERE match_id=$1 ORDER BY id`,
		matchID)
	if err != nil {
		return nil, fmt.Errorf("get attendees: %w", err)
	}
	defer rows.Close()

	var attendees []models.MatchAttendee
	for rows.Next() {
		var a models.MatchAttendee
		if err := rows.Scan(&a.ID, &a.MatchID, &a.PlayerID, &a.GuestCount); err != nil {
			return nil, err
		}
		attendees = append(attendees, a)
	}
	return attendees, rows.Err()
}

func UpsertAttendee(ctx context.Context, pool *pgxpool.Pool, matchID, playerID, guestCount int) (*models.MatchAttendee, error) {
	var a models.MatchAttendee
	err := pool.QueryRow(ctx,
		`INSERT INTO match_attendees (match_id, player_id, guest_count) VALUES ($1,$2,$3)
		 ON CONFLICT (match_id, player_id) DO UPDATE SET guest_count=EXCLUDED.guest_count
		 RETURNING id, match_id, player_id, guest_count`,
		matchID, playerID, guestCount).
		Scan(&a.ID, &a.MatchID, &a.PlayerID, &a.GuestCount)
	if err != nil {
		return nil, fmt.Errorf("upsert attendee: %w", err)
	}
	return &a, nil
}

func UpdateAttendee(ctx context.Context, pool *pgxpool.Pool, matchID, playerID, guestCount int) (*models.MatchAttendee, error) {
	var a models.MatchAttendee
	err := pool.QueryRow(ctx,
		`UPDATE match_attendees SET guest_count=$3 WHERE match_id=$1 AND player_id=$2
		 RETURNING id, match_id, player_id, guest_count`,
		matchID, playerID, guestCount).
		Scan(&a.ID, &a.MatchID, &a.PlayerID, &a.GuestCount)
	if err != nil {
		return nil, fmt.Errorf("update attendee: %w", err)
	}
	return &a, nil
}

func DeleteAttendee(ctx context.Context, pool *pgxpool.Pool, matchID, playerID int) error {
	cmd, err := pool.Exec(ctx,
		`DELETE FROM match_attendees WHERE match_id=$1 AND player_id=$2`, matchID, playerID)
	if err != nil {
		return fmt.Errorf("delete attendee: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("attendee not found")
	}
	return nil
}
