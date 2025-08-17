package controllers

import (
	"fmt"
	"net/http"
	"time"

	"go-backend/database"
	"go-backend/dto"
	"go-backend/models"

	"github.com/gin-gonic/gin"
)

// CreateMessage crée un message entre deux users
func CreateMessage(c *gin.Context) {
	var req dto.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.GlobalError{
			CodeStatus:       http.StatusBadRequest,
			Message:          "Invalid request",
			TechnicalMessage: err.Error(),
		})
		return
	}

	senderID := c.GetString("user_id") // injecté par middleware JWT

	if senderID == req.ReceiverID {
		c.JSON(http.StatusBadRequest, dto.GlobalError{
			CodeStatus:       http.StatusBadRequest,
			Message:          "You cannot send a message to yourself! LOL",
			TechnicalMessage: "Sender and receiver IDs are the same",
		})
		return
	}

	msg := models.Message{
		Content:    req.Content,
		SenderID:   senderID,
		ReceiverID: req.ReceiverID,
		SentAt:     time.Now(),
	}

	if err := database.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.GlobalError{
			CodeStatus:       http.StatusInternalServerError,
			Message:          "Could not save message",
			TechnicalMessage: err.Error(),
		})
		return
	}

	res := dto.MessageResponse{
		ID:         msg.ID,
		Content:    msg.Content,
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		CreatedAt:  msg.SentAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, dto.GlobalSuccess{
		CodeStatus: http.StatusCreated,
		Message:    "Message created",
		Data:       res,
	})
}

func DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	var message models.Message

	if err := database.DB.First(&message, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.GlobalError{
			CodeStatus:       http.StatusNotFound,
			Message:          "Message not found",
			TechnicalMessage: fmt.Sprintf("No message with the given id %s", id),
		})
		return
	}

	if err := database.DB.Delete(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.GlobalError{
			CodeStatus:       http.StatusInternalServerError,
			Message:          "Could not delete message",
			TechnicalMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.GlobalSuccess{
		CodeStatus: http.StatusOK,
		Message:    "Message deleted",
		Data:       nil,
	})
}

func UpdateMessage(c *gin.Context) {
	id := c.Param("id")
	var (
		message    models.Message
		newMessage dto.UpdateMessageRequest
	)

	if err := database.DB.First(&message, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.GlobalError{
			CodeStatus:       http.StatusNotFound,
			Message:          "Message not found",
			TechnicalMessage: fmt.Sprintf("No message with the given id %s", id),
		})
		return
	}

	if err := c.ShouldBindJSON(&newMessage); err != nil {
		c.JSON(http.StatusAccepted, dto.GlobalError{
			CodeStatus:       http.StatusInternalServerError,
			Message:          "Could not update message",
			TechnicalMessage: err.Error(),
		})
		return
	}
	if newMessage.Content == "" || newMessage.Content == message.Content {
		message.ReadAt = time.Now()
		goto end
	}
	message.Content = newMessage.Content

end:
	if err := database.DB.Save(&message).Error; err != nil {
		c.JSON(http.StatusAccepted, dto.GlobalError{
			CodeStatus:       http.StatusInternalServerError,
			Message:          "Could not update message",
			TechnicalMessage: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dto.GlobalSuccess{
		CodeStatus: http.StatusOK,
		Message:    "Message updated",
		Data: dto.GlobalSuccess{
			CodeStatus: http.StatusOK,
			Message:    "Message updated",
			Data:       message,
		},
	})
}

// GetMessages récupère les messages reçus et envoyés par le user connecté
func GetMessages(c *gin.Context) {
	userID := c.GetString("user_id")
	var messages []models.Message

	if err := database.DB.
		Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Order("sent_at desc").
		Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.GlobalError{
			CodeStatus:       http.StatusInternalServerError,
			Message:          "Error loading messages",
			TechnicalMessage: err.Error(),
		})
		return
	}

	var res []dto.MessageResponse
	for _, m := range messages {
		res = append(res, dto.MessageResponse{
			ID:         m.ID,
			Content:    m.Content,
			SenderID:   m.SenderID,
			ReceiverID: m.ReceiverID,
			CreatedAt:  m.SentAt.Format(time.RFC3339),
			ReadAt:     m.ReadAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, dto.GlobalSuccess{
		CodeStatus: http.StatusOK,
		Message:    "Messages retrieved",
		Data:       res,
	})
}
