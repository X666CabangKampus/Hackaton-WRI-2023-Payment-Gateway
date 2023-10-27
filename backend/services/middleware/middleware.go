package middleware

import (
	util "backend-hacktober/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Error struct {
	StatusCode int    `json:"error_code"`
	Message    string `json:"error_message"`
}

type MiddlewareHandlerFunc func(jwtS *util.JWTStruct) gin.HandlerFunc

func MiddlewareJWT(handler MiddlewareHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("JWT")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Error{StatusCode: http.StatusUnauthorized, Message: "Unauthorized"})
			return
		}

		jwtS, err := util.ValidateJWTSign(cookie)
		if err != nil {
			fmt.Println("Err: ", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, Error{StatusCode: http.StatusUnauthorized, Message: "Unauthorized"})
			return
		}

		handler(jwtS)(c)
	}
}
