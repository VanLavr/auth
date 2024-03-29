package delivery

type Response struct {
	Error   string `json:"error"`
	Content any    `json:"content"`
}
