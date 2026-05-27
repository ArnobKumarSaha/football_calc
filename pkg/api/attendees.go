package api

import (
	"encoding/json"
	"net/http"

	"github.com/ArnobKumarSaha/football_calc/pkg/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpsertAttendee(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matchID, err := urlParamInt(r, "id")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid match id")
			return
		}
		var body struct {
			PlayerID   int `json:"player_id"`
			GuestCount int `json:"guest_count"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.PlayerID == 0 {
			writeError(w, http.StatusBadRequest, "player_id required")
			return
		}
		a, err := db.UpsertAttendee(r.Context(), pool, matchID, body.PlayerID, body.GuestCount)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, a)
	}
}

func UpdateAttendee(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matchID, err := urlParamInt(r, "id")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid match id")
			return
		}
		playerID, err := urlParamInt(r, "player_id")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid player id")
			return
		}
		var body struct {
			GuestCount int `json:"guest_count"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "invalid body")
			return
		}
		a, err := db.UpdateAttendee(r.Context(), pool, matchID, playerID, body.GuestCount)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, a)
	}
}

func DeleteAttendee(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matchID, err := urlParamInt(r, "id")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid match id")
			return
		}
		playerID, err := urlParamInt(r, "player_id")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid player id")
			return
		}
		if err := db.DeleteAttendee(r.Context(), pool, matchID, playerID); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
