package delivery

type Response struct {
	Error   string `json:"error"`
	Content any    `json:"content"`
}

type Refresh struct {
	Token string `json:"refresh_token"`
}
