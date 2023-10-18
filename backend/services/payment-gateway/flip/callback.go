package flip

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func AccPaymentCallback(c *gin.Context) {
	data, _ := c.GetPostForm("data")
	token, _ := c.GetPostForm("token")
	if token == ValidationToken {
		fmt.Println("data: " + data)
	} else {
		fmt.Println("Errors with token: " + token)
	}
}
