package models

import "time"

type Player struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type PlayerWithBalance struct {
	Player
	Balance float64 `json:"balance"`
}

type Match struct {
	ID        int        `json:"id"`
	Date      string     `json:"date"`
	TotalBill float64    `json:"total_bill"`
	Notes     *string    `json:"notes,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type MatchAttendee struct {
	ID         int `json:"id"`
	MatchID    int `json:"match_id"`
	PlayerID   int `json:"player_id"`
	GuestCount int `json:"guest_count"`
}

type MatchAttendeeWithDue struct {
	MatchAttendee
	PlayerName string  `json:"player_name"`
	Due        float64 `json:"due"`
}

type MatchDetail struct {
	Match
	Attendees []MatchAttendeeWithDue `json:"attendees"`
	Payments  []Payment              `json:"payments"`
}

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

type BalanceHistory struct {
	PlayerID   int                 `json:"player_id"`
	PlayerName string              `json:"player_name"`
	Matches    []MatchBalanceEntry `json:"matches"`
	Standalone []Payment           `json:"standalone_payments"`
	TotalDue   float64             `json:"total_due"`
	TotalPaid  float64             `json:"total_paid"`
	Balance    float64             `json:"balance"`
}

type MatchBalanceEntry struct {
	MatchID   int     `json:"match_id"`
	MatchDate string  `json:"match_date"`
	Due       float64 `json:"due"`
	Paid      float64 `json:"paid"`
}
