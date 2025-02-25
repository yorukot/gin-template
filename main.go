package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yorukot/go-template/app/routes"

	// _ "github.com/yorukot/go-template/pkg/cache" uncomment this to use cache
	// _ "github.com/yorukot/go-template/pkg/oauth" uncomment this to use oauth

	_ "github.com/yorukot/go-template/pkg/database"
	"github.com/yorukot/go-template/pkg/logger"
	"github.com/yorukot/go-template/pkg/middleware"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	root := gin.New()

	root.SetTrustedProxies([]string{"127.0.0.1"})
	root.StaticFile("/favicon.ico", "./static/favicon.ico")
	root.Use(middleware.CustomLogger())
	root.Use(middleware.ErrorLoggerMiddleware())

	r := root.Group("/api/v" + os.Getenv("VERSION"))

	route(r)

	printAppInfo()

	root.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "resource_not_found",
			"message": "Resource not found",
		})
	})

	if err := root.Run(); err != nil {
		logger.Log.Sugar().Fatal("Server failed to start: %v", err)
	}
}

func printAppInfo() {
	info := fmt.Sprintf(`
	Gin Template API
	Version: %s
	Gin Version: %s
	Domain: %s
	`, os.Getenv("VERSION"), gin.Version, os.Getenv("BASE_URL"))
	logger.Log.Info(info)
}

func route(r *gin.RouterGroup) {
	routes.AuthRoute(r)
	routes.UserRoute(r)
}
