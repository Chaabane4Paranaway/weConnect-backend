// routes.go - stub
package routes

import (
	"go-backend/controllers"
	"go-backend/middlewares"
	"go-backend/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	hub := websocket.NewHub()

	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/verify", controllers.VerifyOTP)
		auth.POST("/login", controllers.Login)
	}

	protected := r.Group("/api", middlewares.AuthMiddleware())
	{
		protected.GET("/me", controllers.Me)

		messaging := protected.Group("/messaging")
		{
			messaging.GET("/messages", controllers.GetMessages)
			messaging.POST("/send", controllers.CreateMessage)
			messaging.DELETE("/delete/:id", controllers.DeleteMessage)
			messaging.PUT("/update/:id", controllers.UpdateMessage)

			r.GET("/ws", websocket.WebSocketHandler(hub))
		}

		friendship := protected.Group("/friendship")
		{
			friendship.POST("/add", controllers.AddFriend)
			friendship.GET("/friends", controllers.GetFriendsHandler(hub))
			// friendship.DELETE("/remove/:id", controllers.RemoveFriend)
			// friendship.GET("/list", controllers.ListFriends)
		}
	}
}
