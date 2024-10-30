package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type BotRepository struct {
	baseUrl string
}

func NewBotRepository(baseUrl string) BotRepository {
	return BotRepository{
		baseUrl: baseUrl,
	}
}

func (botRepository BotRepository) Bid(cardIds []uint16) (int, error) {
	payload := map[string][]uint16{"cards": cardIds}

	body, err := json.Marshal(payload)

	if err != nil {
		return 0, fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/bid", botRepository.baseUrl)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	var result struct {
		Bid int `json:"bid"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Bid, nil
}

func (botRepository BotRepository) Pick(
	handCards []uint16,
	pickableCards []uint16,
	tableCards []uint16,
	observedCards []uint16,
	bid int,
	tricksTaken uint,
	playerIndex int,
	numPlayers int,
) (uint16, error) {
	payload := map[string]interface{}{
		"hand_cards":     handCards,
		"pickable_cards": pickableCards,
		"table_cards":    tableCards,
		"observed_cards": observedCards,
		"bid":            bid,
		"tricks_taken":   tricksTaken,
		"player_index":   playerIndex,
		"num_players":    numPlayers,
	}

	body, err := json.Marshal(payload)

	if err != nil {
		return 0, fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/pick", botRepository.baseUrl)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	var result struct {
		CardId uint16 `json:"card_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.CardId, nil
}
