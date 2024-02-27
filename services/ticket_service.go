package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TicketService struct {
	baseUrl string
}

func NewTicketService(baseUrl string) TicketService {
	return TicketService{baseUrl: baseUrl}
}

func (ticketService TicketService) AcquireUserId(ticketId string) string {
	url := fmt.Sprintf("%s/tickets/acquire", ticketService.baseUrl)
	body := []byte(fmt.Sprintf(`{"ticketId": "%s"}`, ticketId))
	reader := bytes.NewReader(body)

	response, err := http.Post(url, "application/json", reader)

	if err != nil {
		LogService{}.Error(map[string]string{
			"message":  err.Error(),
			"ticketId": ticketId,
			"method":   "TicketService@AcquireUserId",
		})
		return ""
	}

	if response.StatusCode != 200 {
		LogService{}.Error(map[string]string{
			"message":    "Unauthorized",
			"ticketId":   ticketId,
			"method":     "TicketService@AcquireUserId",
			"statusCode": response.Status,
		})
		return ""
	}

	var output struct {
		UserId string `json:"userId"`
	}

	err = json.NewDecoder(response.Body).Decode(&output)

	if err != nil {
		LogService{}.Error(map[string]string{
			"message": err.Error(),
			"method":  "TicketService@AcquireUserId",
		})
		return ""
	}

	return output.UserId
}
