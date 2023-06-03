package responses

type DealResponse struct {
	Round int    `json:"round"`
	Trick int    `json:"trick"`
	Cards []int  `json:"cards"`
	State string `json:"state"`
}

type StartBidding struct {
	EndsAt int64 `json:"endsAt"`
}

type StartPicking struct {
	PlayerId string `json:"playerId"`
	EndsAt   int64  `json:"endsAt"`
	CardIds  []int  `json:"cardIds"`
}

type Pick struct {
	PlayerId string `json:"playerId"`
	CardId   int    `json:"cardId"`
}

type AnnounceTrickWinner struct {
	PlayerId string `json:"playerId"`
	CardId   int    `json:"cardId"`
}
