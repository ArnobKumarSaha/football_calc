package models

type PlayerWithBalance struct {
	Player
	Balance float64 `json:"balance"`
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
