package lottery

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LotteryAPIResponse struct {
	LotteryNumber int `json:"lotteryNumber"`
}

func GetWinningNumber(apiURL string) (int, error) {
	response, err := http.Get(apiURL)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("lottery API returned status: %s", response.Status)
	}

	var data LotteryAPIResponse

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return 0, err
	}

	return data.LotteryNumber, nil
}
