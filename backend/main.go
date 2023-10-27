package main

import (
	"backend-hacktober/services"
	"backend-hacktober/services/email"
	"github.com/gin-gonic/gin"
	"github.com/imrenagi/go-payment/util/localconfig"
)

func main() {
	r := gin.Default()

	config, err := localconfig.LoadConfig("./conf/config.yaml")
	if err != nil {
		panic(err)
	}

	secret, err := localconfig.LoadSecret("./conf/secret.yaml")
	if err != nil {
		panic(err)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello world!"})
	})
	services.NewSrv(r, config, secret).Routes()

	go emailservice.EmailQueueReceiver()
	err = r.Run("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
}
