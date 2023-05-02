package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func IndexTemplate(c *gin.Context) {
	Data := map[string]string{
		"Host": os.Getenv("RABBITMQ_URL"),
	}
	c.HTML(http.StatusOK, "index.html", Data)
}
