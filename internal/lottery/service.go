package lottery

import (
	"celtra-lottery/internal/database"
	"database/sql"
	"fmt"
	"time"
)

type Service struct {
	DB     *sql.DB
	APIURL string
}

func NewService(db *sql.DB, apiURL string) *Service {
	return &Service{
		DB:     db,
		APIURL: apiURL,
	}
}

func (s *Service) EnsureRaffle() error {
	raffle, err := database.GetLatestUnfinishedRaffle(s.DB)
	if err != nil {
		return err
	}

	if raffle == nil {
		now := time.Now()

		_, err := database.CreateRaffle(
			s.DB,
			now,
			now.Add(30*time.Second),
		)

		return err
	}

	if time.Now().Before(raffle.EndsAt) {
		return nil
	}

	winningNumber, err := GetWinningNumber(s.APIURL)
	if err != nil {
		return err
	}

	if err := database.FinishRaffle(
		s.DB,
		raffle.ID,
		winningNumber,
	); err != nil {
		return err
	}

	now := time.Now()

	_, err = database.CreateRaffle(
		s.DB,
		now,
		now.Add(30*time.Second),
	)

	return err
}

func (s *Service) Start() {
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for range ticker.C {
			err := s.EnsureRaffle()
			if err != nil {
				fmt.Println("Lottery service error:", err)
			}
		}
	}()
}

func (s *Service) AddEntry(name string, guess int) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}

	if guess < 1 || guess > 30 {
		return fmt.Errorf("guess must be between 1 and 30")
	}

	raffle, err := database.GetActiveRaffle(s.DB)
	if err != nil {
		return err
	}

	if raffle == nil {
		return fmt.Errorf("no active raffle")
	}

	_, err = database.CreateEntry(
		s.DB,
		raffle.ID,
		name,
		guess,
		time.Now(),
	)

	return err
}
