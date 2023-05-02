package routes

import (
	"rabbit/controllers"

	"github.com/gin-gonic/gin"
)

// InitRoutes initializes the gin engine for routing and server.
func InitRoutes() *gin.Engine {

	app := gin.Default()
	app.GET("/", controllers.IndexTemplate)
	app.GET("/notification/:user", controllers.NotifySingleUser)

	app.GET("/notification/all", controllers.Broadcast)
	return app
}
