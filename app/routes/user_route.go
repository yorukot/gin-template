package routes

import (
	"github.com/gin-gonic/gin"
	userCtrl "github.com/yorukot/go-template/app/controllers/user"
	"github.com/yorukot/go-template/pkg/middleware"
)

func UserRoute(r *gin.RouterGroup) {
	userGroup := r.Group("/user")
	userGroup.Use(middleware.IsAuthorized())

	userGroup.GET("/profile", userCtrl.GetProfile)
}
