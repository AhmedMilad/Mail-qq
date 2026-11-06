package utils

// TODO add attachments
type SendEmailRequest struct {
	Subject string `json:"subject" binding:"required"`
	To      string `json:"to" binding:"required,email"`
	Body    string `json:"body" binding:"required"`
}
