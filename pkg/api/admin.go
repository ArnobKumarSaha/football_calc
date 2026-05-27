package api

import (
	"net/http"

	"github.com/ArnobKumarSaha/football_calc/pkg/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CleanupData(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.CleanupData(r.Context(), pool); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	}
}
