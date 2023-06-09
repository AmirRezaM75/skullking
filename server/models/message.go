package models

// ServerMessage Message from server to client structure
type ServerMessage struct {
	Content    any    `json:"content"`
	Command    string `json:"command"`
	GameId     string `json:"-"`
	ReceiverId string `json:"-"`
}

// ClientMessage Message from client to server structure
type ClientMessage struct {
	Command string `json:"command"`
	Content string `json:"content"`
}
