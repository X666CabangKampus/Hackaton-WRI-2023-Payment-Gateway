package main

import (
	"backend-hacktober/services/gateway"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	gateway.NewSrv(r).Routes()

	err := r.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}
