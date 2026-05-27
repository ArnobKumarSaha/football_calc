package api

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

func ImportData(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		has, err := db.HasAnyData(r.Context(), pool)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if has {
			writeError(w, http.StatusConflict, "database already has data; import is only allowed on an empty database")
			return
		}

		if err := r.ParseMultipartForm(1 << 20); err != nil {
			writeError(w, http.StatusBadRequest, "invalid multipart form")
			return
		}
		f, _, err := r.FormFile("file")
		if err != nil {
			writeError(w, http.StatusBadRequest, "missing file field")
			return
		}
		defer f.Close()

		var entries []db.ImportEntry
		scanner := bufio.NewScanner(f)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) < 2 {
				writeError(w, http.StatusBadRequest, fmt.Sprintf("line %d: expected '<name> <amount>', got %q", lineNum, line))
				return
			}
			name := strings.Join(parts[:len(parts)-1], " ")
			amount, err := strconv.ParseFloat(parts[len(parts)-1], 64)
			if err != nil {
				writeError(w, http.StatusBadRequest, fmt.Sprintf("line %d: invalid amount %q", lineNum, parts[len(parts)-1]))
				return
			}
			entries = append(entries, db.ImportEntry{Name: name, Balance: amount})
		}
		if err := scanner.Err(); err != nil {
			writeError(w, http.StatusBadRequest, "error reading file")
			return
		}
		if len(entries) == 0 {
			writeError(w, http.StatusBadRequest, "file is empty")
			return
		}

		if err := db.ImportInitialData(r.Context(), pool, entries); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "imported": len(entries)})
	}
}
