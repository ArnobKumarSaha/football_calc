package api

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
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
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=matches.csv")
		cw := csv.NewWriter(w)
		cw.Write([]string{"match_id", "date", "total_bill", "player_id", "player_name", "guest_count", "due"})
		for _, row := range rows {
			cw.Write([]string{
				strconv.Itoa(row.MatchID),
				row.Date,
				fmt.Sprintf("%.2f", row.TotalBill),
				strconv.Itoa(row.PlayerID),
				row.PlayerName,
				strconv.Itoa(row.GuestCount),
				fmt.Sprintf("%.2f", row.Due),
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
