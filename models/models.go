package models

import "time"

type Raffle struct {
	ID            int64     `json:"id"`
	WinningNumber *int      `json:"winningNumber"`
	StartedAt     time.Time `json:"startedAt"`
	EndsAt        time.Time `json:"endsAt"`
}

type Entry struct {
	ID        int64     `json:"id"`
	RaffleID  int64     `json:"raffleId"`
	Name      string    `json:"name"`
	Guess     int       `json:"guess"`
	CreatedAt time.Time `json:"createdAt"`
}

// rabs za get
type Result struct {
	RaffleID      int64     `json:"raffleId"`
	WinningNumber int       `json:"winningNumber"`
	Winners       []string  `json:"winners"`
	EndedAt       time.Time `json:"endedAt"`
}
