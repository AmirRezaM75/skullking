package models

type Trick struct {
	number        int
	PickingUserId string
	PickedCards   map[string]CardId
}
