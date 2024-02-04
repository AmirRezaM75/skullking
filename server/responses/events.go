package responses

import "encoding/json"

type PublisherEvent struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func GameCreatedEvent(gameId, lobbyId string) (string, error) {
	type GameCreatedContent struct {
		GameId  string `json:"gameId"`
		LobbyId string `json:"lobbyId"`
	}

	var content = GameCreatedContent{
		GameId:  gameId,
		LobbyId: lobbyId,
	}

	return parseEvent("GameCreated", content)
}

func GameEndedEvent(gameId, lobbyId string) (string, error) {
	type GameEndedContent struct {
		GameId  string `json:"gameId"`
		LobbyId string `json:"lobbyId"`
	}

	var content = GameEndedContent{
		GameId:  gameId,
		LobbyId: lobbyId,
	}

	return parseEvent("GameEnded", content)
}

func parseEvent(eventType string, content any) (string, error) {
	message, err := json.Marshal(content)

	if err != nil {
		return "", err
	}

	event := PublisherEvent{
		Type:    eventType,
		Content: string(message),
	}

	e, err := json.Marshal(event)

	if err != nil {
		return "", err
	}

	return string(e), nil
}
