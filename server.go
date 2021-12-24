package main

import (
	"backend/server"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	/*
	 Connection to the db, handling CORS, setting routes and running the server 
	*/
	server.OpenDB()
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	config.ExposeHeaders = []string{"Location", "x-Created-Id"}
	r.Use(cors.New(config))
	server.SetupRouter(r)

	r.Run(":8080")
	//instance := server.New()

	//instance.Engine.Run()
}
