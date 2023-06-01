package models

type Trick struct {
	Number int
	// This is useful to find out who is eligible when receiving 'PICK' command.
	PickingUserId      string
	PickedCards        map[string]CardId
	WinnerPlayerId     string
	StarterPlayerIndex int
}
