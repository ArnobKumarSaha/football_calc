package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ArnobKumarSaha/football_calc/pkg/db"
	"github.com/ArnobKumarSaha/football_calc/pkg/models"
	"github.com/ArnobKumarSaha/football_calc/pkg/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ListMatches(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 50
		offset := 0
		if v, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && v > 0 {
			limit = v
		}
		if v, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && v >= 0 {
			offset = v
		}
		matches, err := db.ListMatches(r.Context(), pool, limit, offset)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if matches == nil {
			matches = []models.Match{}
		}
		writeJSON(w, http.StatusOK, matches)
	}
}

func CreateMatch(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Date      string   `json:"date"`
			TotalBill float64  `json:"total_bill"`
			Notes     *string  `json:"notes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Date == "" || body.TotalBill <= 0 {
			writeError(w, http.StatusBadRequest, "date and total_bill required")
			return
		}
		m, err := db.CreateMatch(r.Context(), pool, body.Date, body.TotalBill, body.Notes)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, m)
	}
}

func GetMatch(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := urlParamInt(r, "id")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		m, err := db.GetMatch(r.Context(), pool, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		attendees, err := db.GetMatchAttendees(r.Context(), pool, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		svcAttendees := make([]service.AttendeeDue, len(attendees))
		for i, a := range attendees {
			svcAttendees[i] = service.AttendeeDue{
				AttendeeID: a.ID,
				PlayerID:   a.PlayerID,
				GuestCount: a.GuestCount,
			}
		}
		withDues := service.ComputeDues(m.TotalBill, svcAttendees)

		attendeesWithDue := make([]models.MatchAttendeeWithDue, len(withDues))
		playerNames := map[int]string{}
		for _, a := range attendees {
			playerNames[a.PlayerID] = ""
		}
		playerIDList := make([]int, 0, len(playerNames))
		for pid := range playerNames {
			playerIDList = append(playerIDList, pid)
		}
		for _, pid := range playerIDList {
			p, err := db.GetPlayer(r.Context(), pool, pid)
			if err == nil {
				playerNames[pid] = p.Name
			}
		}
		for i, wd := range withDues {
			attendeesWithDue[i] = models.MatchAttendeeWithDue{
				MatchAttendee: models.MatchAttendee{
					ID:         wd.AttendeeID,
					MatchID:    id,
					PlayerID:   wd.PlayerID,
					GuestCount: wd.GuestCount,
				},
				PlayerName: playerNames[wd.PlayerID],
				Due:        wd.Due,
			}
		}

		payments, err := db.ListPayments(r.Context(), pool, nil, &id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if payments == nil {
			payments = []models.Payment{}
		}

		detail := models.MatchDetail{
			Match:     *m,
			Attendees: attendeesWithDue,
			Payments:  payments,
		}
		writeJSON(w, http.StatusOK, detail)
	}
}

func UpdateMatch(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := urlParamInt(r, "id")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		var body struct {
			Date      *string  `json:"date"`
			TotalBill *float64 `json:"total_bill"`
			Notes     *string  `json:"notes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "invalid body")
			return
		}
		m, err := db.UpdateMatch(r.Context(), pool, id, body.Date, body.TotalBill, body.Notes)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, m)
	}
}

func DeleteMatch(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := urlParamInt(r, "id")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if err := db.DeleteMatch(r.Context(), pool, id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
