package controllers

import (
	"go-backend/database"
	"go-backend/dto"
	"go-backend/models"
	"go-backend/websocket"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddFriend(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id missing in context"})
		return
	}
	var req dto.AddFriendshipDTO

	// log.Printf("Received request from user %s", userID)
	// log.Printf("Request details: %+v", req)
	//
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.GlobalError{
			CodeStatus:       http.StatusBadRequest,
			TechnicalMessage: err.Error(),
			Message:          "Invalid JSON",
		})
		log.Println(err)
		return
	}

	var friend models.User
	if err := database.DB.Where("pseudo = ?", req.Pseudo).First(&friend).Error; err != nil {
		c.JSON(404, dto.GlobalError{
			CodeStatus:       http.StatusNotFound,
			TechnicalMessage: err.Error(),
			Message:          "User not found",
		})
		return
	}

	u1 := userID.(string)
	u2 := friend.ID

	if u2 < u1 {
		u1, u2 = u2, u1
	}

	var existing models.Friendship
	err := database.DB.Where("user1_id = ? AND user2_id = ?", u1, u2).First(&existing).Error
	if err == nil {
		// Déjà ami → retourne un message genre "Vous êtes déjà amis"
		c.JSON(http.StatusConflict, dto.GlobalError{
			CodeStatus:       http.StatusConflict,
			TechnicalMessage: "Friendship already exists",
			Message:          "You are already friends",
		})
		return
	}

	newFriend := models.Friendship{User1ID: u1, User2ID: u2}
	if err := database.DB.Create(&newFriend).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.GlobalError{
			CodeStatus:       http.StatusInternalServerError,
			TechnicalMessage: err.Error(),
			Message:          "Failed to create friendship",
		})
		return
	}

	c.JSON(http.StatusCreated, dto.GlobalSuccess{
		CodeStatus: http.StatusCreated,
		Message:    "Friendship added successfully",
		Data:       newFriend,
	})
}

func GetFriendsHandler(hub *websocket.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id missing in context"})
			return
		}
		currentUserID := userID.(string)

		var friendships []models.Friendship
		if err := database.DB.
			Where("user1_id = ? OR user2_id = ?", currentUserID, currentUserID).
			Find(&friendships).Error; err != nil {
			c.JSON(http.StatusInternalServerError, dto.GlobalError{
				CodeStatus:       http.StatusInternalServerError,
				TechnicalMessage: err.Error(),
				Message:          "Failed to fetch friends",
			})
			return
		}

		var friendsList []dto.FriendInfo

		for _, f := range friendships {
			var friendUser models.User
			var friendID string
			if f.User1ID == currentUserID {
				friendID = f.User2ID
			} else {
				friendID = f.User1ID
			}

			if err := database.DB.First(&friendUser, "id = ?", friendID).Error; err != nil {
				continue // skip si user non trouvé
			}

			friendsList = append(friendsList, dto.FriendInfo{
				ID:     friendUser.ID,
				Pseudo: friendUser.Pseudo,
				Online: hub.IsOnline(friendUser.ID), // si tu passes le hub dans le handler
			})
		}

		c.JSON(http.StatusOK, dto.GlobalSuccess{
			CodeStatus: http.StatusOK,
			Data:       toMapList(friendsList),
		})
	}
}

// Helper pour convertir slice en []map[string]interface{} si tu veux réutiliser ton GlobalRes
func toMapList(friends []dto.FriendInfo) []map[string]interface{} {
	var list []map[string]interface{}
	for _, f := range friends {
		list = append(list, map[string]interface{}{
			"id":     f.ID,
			"pseudo": f.Pseudo,
			"online": f.Online,
		})
	}
	return list
}
