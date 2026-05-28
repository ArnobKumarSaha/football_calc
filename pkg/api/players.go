package api

import (
	"encoding/json"
	"net/http"

	"github.com/ArnobKumarSaha/football_calc/pkg/db"
	"github.com/ArnobKumarSaha/football_calc/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ListPlayers(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		players, err := db.AllPlayerBalances(r.Context(), pool)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if players == nil {
			players = []models.PlayerWithBalance{}
		}
		writeJSON(w, http.StatusOK, players)
	}
}

func CreatePlayer(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
			writeError(w, http.StatusBadRequest, "name required")
			return
		}
		p, err := db.CreatePlayer(r.Context(), pool, body.Name)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, p)
	}
}

func DeletePlayer(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := urlParamInt(r, "id")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if err := db.DeletePlayer(r.Context(), pool, id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetPlayerBalance(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := urlParamInt(r, "id")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		h, err := db.GetPlayerBalance(r.Context(), pool, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, h)
	}
}
