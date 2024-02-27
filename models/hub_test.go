package models

import (
	"skullking/constants"
	"skullking/pkg/syncx"
	"testing"
	"time"
)

func TestHub_Cleanup(t *testing.T) {
	games := syncx.Map[string, *Game]{}
	games.Store("786e4150-77ca-11ee-b962-0242ac120002", &Game{
		Id:        "786e4150-77ca-11ee-b962-0242ac120002",
		CreatedAt: time.Now().Add(-15 * time.Minute).Unix(),
		State:     constants.StatePending,
	})
	games.Store("786e481c-77ca-11ee-b962-0242ac120002", &Game{
		Id:        "786e481c-77ca-11ee-b962-0242ac120002",
		CreatedAt: time.Now().Add(-40 * time.Minute).Unix(),
		State:     constants.StatePending,
	})
	games.Store("786e495c-77ca-11ee-b962-0242ac120002", &Game{
		Id:        "786e495c-77ca-11ee-b962-0242ac120002",
		CreatedAt: time.Now().Add(-40 * time.Minute).Unix(),
		State:     constants.StatePicking,
	})
	games.Store("786e4ba0-77ca-11ee-b962-0242ac120002", &Game{
		Id:        "786e4ba0-77ca-11ee-b962-0242ac120002",
		CreatedAt: time.Now().Unix(),
		State:     constants.StatePending,
	})

	hub := Hub{
		Games: games,
	}

	hub.Cleanup()

	if hub.Games.Len() != 3 {
		t.Errorf("Expected 3 games but got %d", hub.Games.Len())
	}

	if _, ok := hub.Games.Load("786e481c-77ca-11ee-b962-0242ac120002"); ok {
		t.Errorf("Expected to remove 786e481c-77ca-11ee-b962-0242ac120002 game.")
	}
}
