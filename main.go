package main

import (
	"log"
	"net/http"

	"celtra-lottery/internal/api"
	"celtra-lottery/internal/database"
	"celtra-lottery/internal/lottery"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.CreateTables(db); err != nil {
		log.Fatal(err)
	}

	apiURL := "https://celtra-lottery-assignment-46e1a1474002.herokuapp.com/api/getLotteryNumber"

	service := lottery.NewService(db, apiURL)

	if err := service.EnsureRaffle(); err != nil {
		log.Fatal(err)
	}

	service.Start()

	log.Println("Lottery service running")

	handler := api.NewHandler(service)

	api.RegisterRoutes(handler)

	log.Println("Server running on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
