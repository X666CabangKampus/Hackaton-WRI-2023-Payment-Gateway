package main

import (
	"backend-hacktober/services"
	"backend-hacktober/services/email"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello world!"})
	})
	services.NewSrv(r).Routes()

	go emailservice.EmailQueueReceiver()
	err := r.Run("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
}
