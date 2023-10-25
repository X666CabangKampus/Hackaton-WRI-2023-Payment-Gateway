package user

import (
	util "backend-hacktober/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Server struct {
	DB *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	return &Server{db}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"Password"`
}

func (s Server) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		err := c.BindJSON(&req)
		if err != nil {
			WriteFailResponse(c, http.StatusBadRequest, Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
			return
		}

		var user User
		s.DB.Where(map[string]any{
			"username": req.Username,
			"password": req.Password,
		}).Find(&user)

		if user.ID == 0 {
			WriteFailResponse(c, http.StatusBadRequest, Error{StatusCode: http.StatusBadRequest, Message: "Username or password is wrong"})
			return
		}

		signedJWT, err := util.CreateJWTSign(&util.JWTStruct{Username: user.Username})
		if err != nil {
			WriteFailResponse(c, http.StatusBadRequest, Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
			return
		}
		c.SetCookie("JWT", signedJWT, 3600*24, "", "", true, true)

		WriteSuccessResponse(c, http.StatusOK, signedJWT, map[string]string{})
	}
}

func MiddlewareJWT(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("JWT")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Error{StatusCode: http.StatusUnauthorized, Message: "Unauthorized"})
			return
		}

		_, err = util.ValidateJWTSign(cookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Error{StatusCode: http.StatusUnauthorized, Message: "Unauthorized"})
			return
		}

		handler(c)
	}
}
