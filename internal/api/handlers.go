package api

import (
	"encoding/json"
	"net/http"

	"celtra-lottery/internal/database"
	"celtra-lottery/internal/lottery"
)

type Handler struct {
	Service *lottery.Service
}

type EntryRequest struct {
	Name  string `json:"name"`
	Guess int    `json:"guess"`
}

func NewHandler(service *lottery.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	var request EntryRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = h.Service.AddEntry(request.Name, request.Guess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Entry created",
	})
}

func (h *Handler) GetState(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodGet {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	raffle, err := database.GetActiveRaffle(h.Service.DB)
	if err != nil {
		http.Error(w, "Failed to get raffle", http.StatusInternalServerError)
		return
	}

	if raffle == nil {
		http.Error(w, "No active raffle", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"raffleId": raffle.ID,
		"endsAt":   raffle.EndsAt,
	})
}
func (h *Handler) GetResults(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodGet {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	results, err := database.GetLatestResults(h.Service.DB)
	if err != nil {
		http.Error(w, "Failed to get results", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(results)
}
