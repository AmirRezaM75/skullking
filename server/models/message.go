package models

// ServerMessage Message from server to client structure
type ServerMessage struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
	Command     string `json:"command"`
	SenderId    string `json:"senderId"`
	GameId      string `json:"-"`
	ReceiverId  string `json:"-"`
}

// ClientMessage Message from client to server structure
type ClientMessage struct {
	Command string `json:"command"`
	Content string `json:"content"`
}
