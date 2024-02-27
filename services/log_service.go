package services

import (
	"fmt"
	"log"
)

type LogService struct {
}

func (logService LogService) Error(payload map[string]string) {
	log.Println("ERROR: " + logService.log(payload))
}

func (logService LogService) Info(payload map[string]string) {
	log.Println("INFO: " + logService.log(payload))
}

func (logService LogService) log(payload map[string]string) string {
	var message string

	for k, v := range payload {
		message += fmt.Sprintf(" %s='%s'", k, v)
	}

	return message
}
