package service

import (
	"math"
)

type AttendeeDue struct {
	AttendeeID int
	PlayerID   int
	GuestCount int
	Due        float64
}

func ComputeDues(totalBill float64, attendees []AttendeeDue) []AttendeeDue {
	totalUnits := 0
	for _, a := range attendees {
		totalUnits += 1 + a.GuestCount
	}
	if totalUnits == 0 {
		return attendees
	}

	perUnit := totalBill / float64(totalUnits)
	lastIdx := len(attendees) - 1
	var sumRounded float64

	result := make([]AttendeeDue, len(attendees))
	copy(result, attendees)

	for i := range lastIdx {
		d := round2(perUnit * float64(1+result[i].GuestCount))
		result[i].Due = d
		sumRounded += d
	}
	result[lastIdx].Due = round2(totalBill - sumRounded)
	return result
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
