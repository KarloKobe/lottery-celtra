package database

import (
	"celtra-lottery/models"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Connected to database")

	return db, nil
}

func CreateTables(db *sql.DB) error {
	rafflesTable := `
	CREATE TABLE IF NOT EXISTS raffles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		winning_number INTEGER,
		started_at DATETIME NOT NULL,
		ends_at DATETIME NOT NULL
	);
	`

	entriesTable := `
	CREATE TABLE IF NOT EXISTS entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		raffle_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		guess INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (raffle_id) REFERENCES raffles(id)
	);
	`

	if _, err := db.Exec(rafflesTable); err != nil {
		return err
	}

	if _, err := db.Exec(entriesTable); err != nil {
		return err
	}

	fmt.Println("Database tables ready")

	return nil
}

func CreateRaffle(db *sql.DB, startedAt time.Time, endsAt time.Time) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO raffles (started_at, ends_at)
		VALUES (?, ?)
	`, startedAt, endsAt)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetActiveRaffle(db *sql.DB) (*models.Raffle, error) {
	row := db.QueryRow(`
		SELECT id, winning_number, started_at, ends_at
		FROM raffles
		WHERE winning_number IS NULL
		AND ends_at > ?
		ORDER BY id DESC
		LIMIT 1
	`, time.Now())

	var raffle models.Raffle

	err := row.Scan(
		&raffle.ID,
		&raffle.WinningNumber,
		&raffle.StartedAt,
		&raffle.EndsAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &raffle, nil
}

func FinishRaffle(db *sql.DB, raffleID int64, winningNumber int) error {
	_, err := db.Exec(`
		UPDATE raffles
		SET winning_number = ?
		WHERE id = ?
	`, winningNumber, raffleID)

	if err != nil {
		return err
	}

	return nil
}

func GetLatestUnfinishedRaffle(db *sql.DB) (*models.Raffle, error) {
	row := db.QueryRow(`
		SELECT id, winning_number, started_at, ends_at
		FROM raffles
		WHERE winning_number IS NULL
		ORDER BY id DESC
		LIMIT 1
	`)

	var raffle models.Raffle

	err := row.Scan(
		&raffle.ID,
		&raffle.WinningNumber,
		&raffle.StartedAt,
		&raffle.EndsAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &raffle, nil
}

func CreateEntry(
	db *sql.DB,
	raffleID int64,
	name string,
	guess int,
	createdAt time.Time,
) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO entries (raffle_id, name, guess, created_at)
		VALUES (?, ?, ?, ?)
	`, raffleID, name, guess, createdAt)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetLatestResults(db *sql.DB) ([]models.Result, error) {
	rows, err := db.Query(`
		SELECT id, winning_number, ends_at
		FROM raffles
		WHERE winning_number IS NOT NULL
		ORDER BY id DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Result

	for rows.Next() {
		var result models.Result

		err := rows.Scan(
			&result.RaffleID,
			&result.WinningNumber,
			&result.EndedAt,
		)
		if err != nil {
			return nil, err
		}

		winnerRows, err := db.Query(`
			SELECT name
			FROM entries
			WHERE raffle_id = ?
			AND guess = ?
		`, result.RaffleID, result.WinningNumber)

		if err != nil {
			return nil, err
		}

		for winnerRows.Next() {
			var name string

			if err := winnerRows.Scan(&name); err != nil {
				winnerRows.Close()
				return nil, err
			}

			result.Winners = append(result.Winners, name)
		}

		winnerRows.Close()

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
