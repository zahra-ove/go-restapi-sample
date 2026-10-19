package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zahra-ove/go-restapi-sample/db"
	"github.com/zahra-ove/go-restapi-sample/routes"
)

func main() {
	// @TODO: air does not hot reload. check it why
	// @TODO: how to use delv for debugging inside container

	db.InitDB()

	server := gin.Default()
	routes.RegisterRoutes(server)

	server.Run(":8080")
}
