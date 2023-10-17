package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const ValidationToken = "$2y$13$e2TMXMauN6U0fzjNyJJE2ufTXr16/iUF9LKQjAaZDp4D3gybtXtUa"

func CallbackHooks(c *gin.Context) {
	data, _ := c.GetPostForm("data")
	token, _ := c.GetPostForm("token")
	if token == ValidationToken {
		fmt.Println("data: " + data)
	} else {
		fmt.Println("Errors on token with token: " + token)
	}
}
