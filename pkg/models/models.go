package models

import "time"

// Player maps to the players table.
type Player struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// Match maps to the matches table; TotalBill is split among attendees.
type Match struct {
	ID        int        `json:"id"`
	Date      string     `json:"date"`
	TotalBill float64    `json:"total_bill"`
	Notes     *string    `json:"notes,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// MatchAttendee maps to the match_attendees table; GuestCount adds extra heads to the due calc.
type MatchAttendee struct {
	ID         int `json:"id"`
	MatchID    int `json:"match_id"`
	PlayerID   int `json:"player_id"`
	GuestCount int `json:"guest_count"`
}

// Payment maps to the payments table; MatchID is nil for standalone (non-match) payments.
type Payment struct {
	ID         int       `json:"id"`
	PlayerID   int       `json:"player_id"`
	PlayerName string    `json:"player_name"`
	MatchID    *int      `json:"match_id,omitempty"`
	Amount     float64   `json:"amount"`
	Notes      *string   `json:"notes,omitempty"`
	PaidAt     string    `json:"paid_at"`
	CreatedAt  time.Time `json:"created_at"`
}
