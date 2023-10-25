package main

import (
	"backend-hacktober/services"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	services.NewSrv(r).Routes()

	err := r.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}
