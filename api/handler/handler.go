package handler

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"rabbitmq/internal/models"
	"rabbitmq/internal/rabbitmq"

	"github.com/gin-gonic/gin"
)

func SendHandler(c *gin.Context) {
	var body []byte
	var user models.User
	var contentType string
	var err error

	if body, err = io.ReadAll(c.Request.Body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}

	contentType = c.GetHeader("Content-type")

	if contentType == "application/json"{
		if err = json.Unmarshal(body, &user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode JSON"})
			return
		}
	}else if contentType == "application/xml"{
		if  err = xml.Unmarshal(body, &user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode XML"})
			return
		}
	}else{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported Content-Type"})
		return
	}
	
	rabbitmq.SendToQueue(body, contentType)
	c.String(http.StatusOK, "Message sent to RabbitMQ\n")
}
