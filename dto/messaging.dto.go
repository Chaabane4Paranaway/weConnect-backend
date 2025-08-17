package dto

type CreateMessageRequest struct {
	ReceiverID string `json:"recipient_id" binding:"required,uuid"`
	Content    string `json:"content" binding:"required"`
}

type MessageResponse struct {
	ID         string `json:"id"`
	Content    string `json:"content"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"recipient_id"`
	CreatedAt  string `json:"created_at"`
	ReadAt     string `json:"read_at,omitempty"`
}

type UpdateMessageRequest struct {
	Content string `json:"content"`
}
