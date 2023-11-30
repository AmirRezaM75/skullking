package responses

// @link https://github.com/golang/go/issues/44692
// I can't use uint8 for cardId. This is due to byte being an alias for uint8

type Deal struct {
	Round int      `json:"round"`
	Trick int      `json:"trick"`
	Cards []uint16 `json:"cards"`
	State string   `json:"state"`
}

// When the game begins, we compute a unique index for each player
// and need to inform all players about sequence.
type Started struct {
	State   string   `json:"state"`
	Players []Player `json:"players"`
}

type StartBidding struct {
	EndsAt int64  `json:"endsAt"`
	State  string `json:"state"`
	Round  int    `json:"round"`
}

type StartPicking struct {
	PlayerId string   `json:"playerId"`
	EndsAt   int64    `json:"endsAt"`
	CardIds  []uint16 `json:"cardIds"`
	State    string   `json:"state"`
}

type Picked struct {
	PlayerId string `json:"playerId"`
	CardId   uint16 `json:"cardId"`
}

type AnnounceTrickWinner struct {
	PlayerId string `json:"playerId"`
	CardId   uint16 `json:"cardId"`
}

type Card struct {
	Id     uint16 `json:"id"`
	Number int    `json:"number"`
	Type   string `json:"type"`
}

type Error struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type EndBidding struct {
	Bids []Bid `json:"bids"`
}

type Bid struct {
	PlayerId string `json:"playerId"`
	Number   int    `json:"number"`
}

// Bade is rarely used in the context of inviting or requesting someone to do something,
// But I prefer to make this grammatical mistake to distinguish my commands.
type Bade struct {
	Number int `json:"number"`
}

type NextTrick struct {
	Round int `json:"round"`
	Trick int `json:"trick"`
}

type AnnounceScore struct {
	Scores []Score `json:"scores"`
}

type Score struct {
	PlayerId string `json:"playerId"`
	Score    int    `json:"score"`
}

type CreateGame struct {
	Id string `json:"id"`
}

type Player struct {
	Id              string   `json:"id"`
	Username        string   `json:"username"`
	Avatar          string   `json:"avatar"`
	Score           int      `json:"score"`
	Bid             int      `json:"bid"`
	HandCardIds     []uint16 `json:"handCardIds"`
	PickableCardIds []uint16 `json:"pickableCardIds"`
	WonTricksCount  uint     `json:"wonTricksCount"`
}

type TableCard struct {
	PlayerId string `json:"playerId"`
	CardId   uint16 `json:"cardId"`
}

type Init struct {
	Round          int         `json:"round"`
	Trick          int         `json:"trick"`
	State          string      `json:"state"`
	ExpirationTime int64       `json:"expirationTime"`
	PickingUserId  string      `json:"pickingUserId"`
	Players        []Player    `json:"players"`
	CreatorId      string      `json:"creatorId"`
	TableCards     []TableCard `json:"tableCards"`
}

type Left struct {
	PlayerId string `json:"playerId"`
}
