package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zahra-ove/go-restapi-sample/db"
	"github.com/zahra-ove/go-restapi-sample/routes"
)

func main() {
	db.InitDB()

	server := gin.Default()
	routes.RegisterRoutes(server)

	server.Run(":8080")
}
