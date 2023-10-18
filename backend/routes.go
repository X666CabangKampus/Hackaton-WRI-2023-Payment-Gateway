package main

import (
	"backend-hacktober/services/payment-gateway/flip"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"test": "success",
		})
	})

	v1 := r.Group("/v1")
	{
		rFlip := v1.Group("/flip")
		{
			rFlip.GET("/callback")
			rFlip.POST("/callback", flip.AccPaymentCallback)
		}
	}

	fmt.Println("RUNNING SERVER...")
	err := r.Run()
	if err != nil {
		panic(err)
		return
	}
}
