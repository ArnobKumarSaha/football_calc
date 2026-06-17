package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// exportTable pairs a header comment with a query that emits one fully-formed
// INSERT statement (as text) per row. Quoting/escaping is delegated to
// PostgreSQL's quote_nullable(), which handles NULLs, embedded quotes, and
// timestamp/numeric casts correctly.
type exportTable struct {
	name  string
	query string
}

// Order matters: parents before children so FK references resolve on replay.
var exportTables = []exportTable{
	{
		name: "players",
		query: `SELECT 'INSERT INTO players (id, name, created_at, deleted_at) VALUES (' ||
			id || ', ' || quote_nullable(name) || ', ' || quote_nullable(created_at) || ', ' ||
			quote_nullable(deleted_at) || ');'
			FROM players ORDER BY id`,
	},
	{
		name: "matches",
		query: `SELECT 'INSERT INTO matches (id, date, total_bill, notes, created_at, deleted_at) VALUES (' ||
			id || ', ' || quote_nullable(date) || ', ' || quote_nullable(total_bill) || ', ' ||
			quote_nullable(notes) || ', ' || quote_nullable(created_at) || ', ' ||
			quote_nullable(deleted_at) || ');'
			FROM matches ORDER BY id`,
	},
	{
		name: "match_attendees",
		query: `SELECT 'INSERT INTO match_attendees (id, match_id, player_id, guest_count) VALUES (' ||
			id || ', ' || match_id || ', ' || player_id || ', ' || guest_count || ');'
			FROM match_attendees ORDER BY id`,
	},
	{
		name: "payments",
		query: `SELECT 'INSERT INTO payments (id, player_id, match_id, amount, notes, paid_at, created_at) VALUES (' ||
			id || ', ' || player_id || ', ' || quote_nullable(match_id) || ', ' || quote_nullable(amount) || ', ' ||
			quote_nullable(notes) || ', ' || quote_nullable(paid_at) || ', ' || quote_nullable(created_at) || ');'
			FROM payments ORDER BY id`,
	},
}

// ExportSQL produces a self-contained SQL script (schema + all data + sequence
// resets) that re-initializes a fresh, empty PostgreSQL with an exact clone of
// the current data, including soft-deleted rows. Suitable as a KubeDB
// initialization script_source.
func ExportSQL(ctx context.Context, pool *pgxpool.Pool) (string, error) {
	var b strings.Builder

	fmt.Fprintf(&b, "-- football-calc full backup — generated %s\n", time.Now().UTC().Format(time.RFC3339))
	b.WriteString("BEGIN;\n\n")

	// Schema first so the script runs on a truly empty database.
	b.WriteString(SchemaSQL())
	b.WriteString("\n")

	for _, t := range exportTables {
		rows, err := pool.Query(ctx, t.query)
		if err != nil {
			return "", fmt.Errorf("export %s: %w", t.name, err)
		}
		fmt.Fprintf(&b, "\n-- %s\n", t.name)
		for rows.Next() {
			var stmt string
			if err := rows.Scan(&stmt); err != nil {
				rows.Close()
				return "", fmt.Errorf("scan %s: %w", t.name, err)
			}
			b.WriteString(stmt)
			b.WriteString("\n")
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return "", fmt.Errorf("read %s: %w", t.name, err)
		}
		rows.Close()
	}

	// Reset SERIAL sequences so future inserts continue past the max id.
	b.WriteString("\n-- reset sequences so future inserts don't collide with explicit ids\n")
	// For an empty table, MAX(id) is NULL: set the sequence to 1 with
	// is_called=false so the first nextval() returns 1 (setval rejects 0).
	for _, t := range exportTables {
		fmt.Fprintf(&b, "SELECT setval('%s_id_seq', GREATEST((SELECT COALESCE(MAX(id),1) FROM %s), 1), (SELECT MAX(id) FROM %s) IS NOT NULL);\n", t.name, t.name, t.name)
	}

	b.WriteString("\nCOMMIT;\n")
	return b.String(), nil
}
