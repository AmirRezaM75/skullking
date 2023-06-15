package models

import (
	"github.com/AmirRezaM75/skull-king/constants"
	"testing"
	"time"
)

func TestHub_Cleanup(t *testing.T) {
	hub := Hub{
		Games: map[string]*Game{
			"Legolas":   {Id: "Legolas", CreatedAt: time.Now().Add(-15 * time.Minute).Unix(), State: constants.StatePending},
			"Galadriel": {Id: "Galadriel", CreatedAt: time.Now().Add(-40 * time.Minute).Unix(), State: constants.StatePending},
			"Sauron":    {Id: "Sauron", CreatedAt: time.Now().Add(-40 * time.Minute).Unix(), State: constants.StatePicking},
			"Elrond":    {Id: "Elrond", CreatedAt: time.Now().Unix(), State: constants.StatePending},
		},
	}

	hub.Cleanup()

	if len(hub.Games) != 3 {
		t.Errorf("Expected 3 games but got %d", len(hub.Games))
	}

	if _, ok := hub.Games["Galadriel"]; ok {
		t.Errorf("Expected to remove Galadriel game.")
	}
}
