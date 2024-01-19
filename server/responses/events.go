package responses

import "encoding/json"

type PublisherEvent struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func NewGameCreatedEvent(gameId, lobbyId string) (string, error) {
	type GameCreatedContent struct {
		GameId  string `json:"gameId"`
		LobbyId string `json:"lobbyId"`
	}

	var content = GameCreatedContent{
		GameId:  gameId,
		LobbyId: lobbyId,
	}

	message, err := json.Marshal(content)

	if err != nil {
		return "", err
	}

	event := PublisherEvent{
		Type:    "GameCreated",
		Content: string(message),
	}

	e, err := json.Marshal(event)

	if err != nil {
		return "", err
	}

	return string(e), nil
}
