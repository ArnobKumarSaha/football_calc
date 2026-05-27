package api

import (
	"encoding/json"
	"net/http"

	"github.com/ArnobKumarSaha/football_calc/pkg/db"
	"github.com/ArnobKumarSaha/football_calc/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ListPayments(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playerID := queryParamInt(r, "player_id")
		matchID := queryParamInt(r, "match_id")
		payments, err := db.ListPayments(r.Context(), pool, playerID, matchID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if payments == nil {
			payments = []models.Payment{}
		}
		writeJSON(w, http.StatusOK, payments)
	}
}

func CreatePayment(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			PlayerID int     `json:"player_id"`
			MatchID  *int    `json:"match_id"`
			Amount   float64 `json:"amount"`
			PaidAt   string  `json:"paid_at"`
			Notes    *string `json:"notes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.PlayerID == 0 || body.Amount <= 0 || body.PaidAt == "" {
			writeError(w, http.StatusBadRequest, "player_id, amount, paid_at required")
			return
		}
		p, err := db.CreatePayment(r.Context(), pool, body.PlayerID, body.MatchID, body.Amount, body.PaidAt, body.Notes)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, p)
	}
}

func UpdatePayment(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := urlParamInt(r, "id")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		var body struct {
			MatchID *int     `json:"match_id"`
			Amount  *float64 `json:"amount"`
			PaidAt  *string  `json:"paid_at"`
			Notes   *string  `json:"notes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "invalid body")
			return
		}
		p, err := db.UpdatePayment(r.Context(), pool, id, body.MatchID, body.Amount, body.PaidAt, body.Notes)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, p)
	}
}

func DeletePayment(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := urlParamInt(r, "id")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if err := db.DeletePayment(r.Context(), pool, id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
