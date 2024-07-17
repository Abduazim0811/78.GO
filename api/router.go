package api

import (
	"rabbitmq/api/handler"

	"github.com/gin-gonic/gin"
)


func Router(){
	router := gin.Default()

	router.POST("/send", handler.SendHandler)

	router.Run(":7777")

}