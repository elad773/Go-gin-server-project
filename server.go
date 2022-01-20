package main

import (
	"backend/handlers"
	"backend/model"
	"backend/server"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	/*
	 Connection to the db, handling CORS, setting routes and running the server
	*/
	model.OpenDB()
	r := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("notemptystring", handlers.NotEmptyString)
	}

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	config.ExposeHeaders = []string{"Location", "x-Created-Id"}
	r.Use(cors.New(config))
	server.SetupRouter(r)

	r.Run(":8080")

}
