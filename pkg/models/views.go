package models

// PlayerWithBalance is returned by GET /api/players; Balance = SUM(payments) - SUM(dues).
type PlayerWithBalance struct {
	Player
	Balance float64 `json:"balance"`
}

// MatchAttendeeWithDue enriches MatchAttendee with the player name and computed due for that match.
type MatchAttendeeWithDue struct {
	MatchAttendee
	PlayerName string  `json:"player_name"`
	Due        float64 `json:"due"`
}

// MatchDetail is the full response for GET /api/matches/{id}: match + per-attendee dues + all payments.
type MatchDetail struct {
	Match
	Attendees []MatchAttendeeWithDue `json:"attendees"`
	Payments  []Payment              `json:"payments"`
}

// BalanceHistory is returned by GET /api/players/{id}/balance; breaks down dues and payments per match.
type BalanceHistory struct {
	PlayerID   int                 `json:"player_id"`
	PlayerName string              `json:"player_name"`
	Matches    []MatchBalanceEntry `json:"matches"`
	Standalone []Payment           `json:"standalone_payments"`
	TotalDue   float64             `json:"total_due"`
	TotalPaid  float64             `json:"total_paid"`
	Balance    float64             `json:"balance"`
}

// MatchBalanceEntry is one row inside BalanceHistory.Matches: what the player owed vs paid for a match.
type MatchBalanceEntry struct {
	MatchID   int     `json:"match_id"`
	MatchDate string  `json:"match_date"`
	Due       float64 `json:"due"`
	Paid      float64 `json:"paid"`
}
