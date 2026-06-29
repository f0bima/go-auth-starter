package main

import (
	openapiui "github.com/PeterTakahashi/gin-openapi/openapiui"
	"github.com/f0bima/go-auth-starter/internal/feature"
	"github.com/f0bima/go-core/bootstrap"
	"github.com/gin-gonic/gin"
)

// @title           Auth API
// @version         1.0
// @description     This is an authentication service API.
// @host      localhost:8080
// @BasePath  /
func main() {
	// Initialize core foundation (config, logger, telemetry, database)
	app := bootstrap.Bootstrap("../../.env")

	// Setup router with middleware
	app.Router = gin.New()
	app.Router.Use(gin.Recovery())
	app.Router.Use(gin.Logger())

	// Register auth module (repository, usecase, controller, routes)
	feature.Register(app)

	// Serve API docs at /docs
	app.Router.GET("/docs/*any", openapiui.WrapHandler(openapiui.Config{
		SpecURL:      "/docs/openapi.json",
		SpecFilePath: "./docs/swagger.json",
		Title:        "Auth API Docs",
		Theme:        "light",
	}))

	// Ping endpoint
	app.Router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Start server with graceful shutdown
	app.Run()
}
