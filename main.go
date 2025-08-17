// @title GO_Backend
// @version 1.1
// @description Swagger pour apis gobackend
// @host localhost:8080
// @BasePath /

package main

import (
	db "go-backend/database"
	_ "go-backend/docs"
	"go-backend/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	db.InitDB()
	// DB := db.DB

	r := gin.Default()
	routes.SetupRoutes(r)
	// test()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}

func test() {

}
