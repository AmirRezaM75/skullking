package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LobbyService struct {
	baseUrl string
	token   string
}

func NewLobbyService(baseUrl, token string) LobbyService {
	return LobbyService{baseUrl: baseUrl, token: token}
}

type LobbyResponse struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Players   []Player `json:"players"`
	Bots      []Bot    `json:"bots"`
	CreatorId string   `json:"creatorId"`
	CreatedAt int64    `json:"createdAt"`
}

type Bot struct {
	Id       uint8  `json:"id"`
	Username string `json:"username"`
	AvatarId uint8  `json:"avatarId"`
}

type Player struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	AvatarId uint8  `json:"avatarId"`
}

func (lobbyService LobbyService) FindById(id string) *LobbyResponse {
	url := fmt.Sprintf("%s/lobbies/%s", lobbyService.baseUrl, id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Token", lobbyService.token)

	if err != nil {
		LogService{}.Error(map[string]string{
			"url":     url,
			"message": err.Error(),
			"method":  "LobbyService@FindById",
		})
		return nil
	}

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		LogService{}.Error(map[string]string{
			"url":     url,
			"message": err.Error(),
			"method":  "LobbyService@FindById",
		})
		return nil
	}

	if response.StatusCode != 200 {
		LogService{}.Error(map[string]string{
			"method":     "LobbyService@FindById",
			"statusCode": response.Status,
		})
		return nil
	}

	var lobby LobbyResponse

	err = json.NewDecoder(response.Body).Decode(&lobby)

	if err != nil {
		LogService{}.Error(map[string]string{
			"message":     err.Error(),
			"description": "Could not decode LobbyResponse.",
			"method":      "LobbyService@FindById",
		})
		return nil
	}

	return &lobby
}
