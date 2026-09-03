package models

import "time"

type Raffle struct {
	ID            int64
	WinningNumber *int
	StartedAt     time.Time
	EndsAt        time.Time
}

type Entry struct {
	ID        int64
	RaffleID  int64
	Name      string
	Guess     int
	CreatedAt time.Time
}

// rabs za get
type Result struct {
	RaffleID      int64
	WinningNumber int
	Winners       []string
	EndedAt       time.Time
}
