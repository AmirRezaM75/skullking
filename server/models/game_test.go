package models

import (
	"skullking/pkg/syncx"
	"testing"
)

func TestGame_GetAvailableAvatar(t *testing.T) {
	var players syncx.Map[string, *Player]

	players.Store("Hermione", &Player{
		Avatar: "1.jpg",
	})
	players.Store("Harry", &Player{
		Avatar: "3.jpg",
	})
	players.Store("Dobby", &Player{
		Avatar: "4.jpg",
	})

	game := Game{
		Players: players,
	}

	image := game.GetAvailableAvatar()

	if image != "2.jpg" {
		t.Errorf("Expected 2.jpg as avatar, got %s", image)
	}
}
