package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init(logger *log.Logger) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logger.Println("WARN: DATABASE_URL not set, running without DB")
		return
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatalf("ERROR: Failed to connect to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		logger.Fatalf("ERROR: Database ping failed: %v", err)
	}

	runMigrations(logger)
	logger.Println("INFO: Database connected and migrated")
}

func runMigrations(logger *log.Logger) {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS leads (
			id                SERIAL PRIMARY KEY,
			full_name         TEXT NOT NULL,
			business_name     TEXT NOT NULL,
			email             TEXT NOT NULL,
			phone             TEXT NOT NULL,
			website           TEXT,
			message           TEXT NOT NULL,
			source            TEXT NOT NULL DEFAULT 'contact_form',
			lead_score        INTEGER,
			status            TEXT NOT NULL DEFAULT 'new',
			stripe_session_id TEXT,
			created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS partner_applications (
			id                SERIAL PRIMARY KEY,
			company_name      TEXT NOT NULL,
			contact_name      TEXT NOT NULL,
			contact_email     TEXT NOT NULL,
			contact_phone     TEXT,
			website           TEXT,
			years_in_business TEXT,
			client_count      TEXT,
			why_partner       TEXT NOT NULL,
			expected_volume   TEXT,
			status            TEXT NOT NULL DEFAULT 'pending',
			internal_notes    TEXT,
			created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
	}

	for _, m := range migrations {
		if _, err := DB.Exec(m); err != nil {
			logger.Fatalf("ERROR: Migration failed: %v", err)
		}
	}
}
