package ws

type Game struct {
	id             string
	round          int
	trick          int
	state          string
	expirationTime int
	players        map[int]*Player
	rounds         map[int]*Round
}

type Round struct {
	number         int
	scores         map[int]int
	dealtCards     map[int][]CardId
	remainingCards map[int][]CardId
	bids           map[int]int
	tricks         []Trick
}

type Trick struct {
	number        int
	pickingUserId int
	pickedCards   map[int]CardId
}

func getNextPlayerIdForPicking(game Game, trick Trick) int {
	var currentPickingPlayerHaveFound = false

	var pickerId int

	for playerId, _ := range game.players {
		if currentPickingPlayerHaveFound {
			pickerId = playerId
		}

		if playerId == trick.pickingUserId {
			currentPickingPlayerHaveFound = true
		}
	}

	return pickerId
}
