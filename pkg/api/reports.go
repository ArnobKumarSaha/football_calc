package api

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ArnobKumarSaha/football_calc/pkg/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MonthlyReport(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		year := time.Now().Year()
		if v, err := strconv.Atoi(r.URL.Query().Get("year")); err == nil {
			year = v
		}
		reports, err := db.MonthlyReports(r.Context(), pool, year)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if reports == nil {
			reports = []db.MonthlyReport{}
		}
		writeJSON(w, http.StatusOK, reports)
	}
}

func AttendanceReport(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reports, err := db.AttendanceReports(r.Context(), pool)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if reports == nil {
			reports = []db.AttendanceReport{}
		}
		writeJSON(w, http.StatusOK, reports)
	}
}

func MatchesCSV(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.MatchesCSV(r.Context(), pool)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		type matchRow struct {
			date      string
			totalBill float64
			players   []string
		}
		order := []int{}
		matches := map[int]*matchRow{}
		for _, row := range rows {
			if _, ok := matches[row.MatchID]; !ok {
				matches[row.MatchID] = &matchRow{date: row.Date, totalBill: row.TotalBill}
				order = append(order, row.MatchID)
			}
			name := row.PlayerName
			if row.GuestCount > 0 {
				name += fmt.Sprintf("+%d", row.GuestCount)
			}
			matches[row.MatchID].players = append(matches[row.MatchID].players, name)
		}

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=matches.csv")
		cw := csv.NewWriter(w)
		cw.Write([]string{"match_id", "date", "total_bill", "players"})
		for _, id := range order {
			m := matches[id]
			cw.Write([]string{
				strconv.Itoa(id),
				m.date,
				fmt.Sprintf("%.2f", m.totalBill),
				strings.Join(m.players, ", "),
			})
		}
		cw.Flush()
	}
}

func PaymentsCSV(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.PaymentsCSV(r.Context(), pool)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=payments.csv")
		cw := csv.NewWriter(w)
		cw.Write([]string{"id", "player_id", "player_name", "match_id", "amount", "paid_at", "notes"})
		for _, row := range rows {
			matchID := ""
			if row.MatchID != nil {
				matchID = strconv.Itoa(*row.MatchID)
			}
			notes := ""
			if row.Notes != nil {
				notes = *row.Notes
			}
			cw.Write([]string{
				strconv.Itoa(row.ID),
				strconv.Itoa(row.PlayerID),
				row.PlayerName,
				matchID,
				fmt.Sprintf("%.2f", row.Amount),
				row.PaidAt,
				notes,
			})
		}
		cw.Flush()
	}
}
