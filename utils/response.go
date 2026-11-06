package utils

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type EmailLogResponse struct {
	ID        uint   `json:"id"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Body      string `json:"body,omitempty"`
	Subject   string `json"subject"`
}
