package flip

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func CallbackHooks(c *gin.Context) {
	data, _ := c.GetPostForm("data")
	token, _ := c.GetPostForm("token")
	if token == ValidationToken {
		fmt.Println("data: " + data)
	} else {
		fmt.Println("Errors on token with token: " + token)
	}
}
