package api

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(pool *pgxpool.Pool, adminPassword string, static fs.FS) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	r.Get("/readyz", func(w http.ResponseWriter, req *http.Request) {
		if err := pool.Ping(req.Context()); err != nil {
			writeError(w, http.StatusServiceUnavailable, "db unavailable")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Get("/api/players", ListPlayers(pool))
	r.Get("/api/players/{id}/balance", GetPlayerBalance(pool))

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware(adminPassword))
		r.Post("/api/players", CreatePlayer(pool))
		r.Delete("/api/players/{id}", DeletePlayer(pool))
	})

	r.Get("/api/matches", ListMatches(pool))
	r.Get("/api/matches/{id}", GetMatch(pool))

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware(adminPassword))
		r.Post("/api/matches", CreateMatch(pool))
		r.Patch("/api/matches/{id}", UpdateMatch(pool))
		r.Delete("/api/matches/{id}", DeleteMatch(pool))

		r.Post("/api/matches/{id}/attendees", UpsertAttendee(pool))
		r.Patch("/api/matches/{id}/attendees/{player_id}", UpdateAttendee(pool))
		r.Delete("/api/matches/{id}/attendees/{player_id}", DeleteAttendee(pool))
	})

	r.Get("/api/payments", ListPayments(pool))

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware(adminPassword))
		r.Post("/api/payments", CreatePayment(pool))
		r.Patch("/api/payments/{id}", UpdatePayment(pool))
		r.Delete("/api/payments/{id}", DeletePayment(pool))
	})

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware(adminPassword))
		r.Get("/api/auth", func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
		})
		r.Post("/api/admin/cleanup", CleanupData(pool))
		r.Post("/api/admin/import", ImportData(pool))
	})

	r.Get("/api/reports/monthly", MonthlyReport(pool))
	r.Get("/api/reports/matches.csv", MatchesCSV(pool))
	r.Get("/api/reports/payments.csv", PaymentsCSV(pool))

	r.Get("/api/worldcup/fixtures", WorldCupFixtures())
	r.Get("/api/worldcup/standings", WorldCupStandings())

	r.Handle("/*", http.FileServer(http.FS(static)))

	return r
}

func authMiddleware(password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token != "Bearer "+password {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
