package user

import (
	"backend-hacktober/services/middleware"
	"backend-hacktober/services/user/manage"
	userManage "backend-hacktober/services/user/manage"
	"backend-hacktober/services/user/model"
	util "backend-hacktober/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	Manager userManage.Manager
}

func NewServer(manager userManage.Manager) *Server {
	return &Server{manager}
}

func (S Server) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req manage.LoginRequest
		err := c.BindJSON(&req)
		if err != nil {
			util.WriteFailResponse(c, http.StatusBadRequest, util.Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
			return
		}

		user, err := S.Manager.Login(req)
		if err != nil || user == nil {
			util.WriteFailResponse(c, http.StatusBadRequest, util.Error{StatusCode: http.StatusBadRequest, Message: "Username or password is wrong"})
			return
		}

		signedJWT, err := util.CreateJWTSign(&util.JWTStruct{Username: user.Username})
		if err != nil {
			util.WriteFailResponse(c, http.StatusBadRequest, util.Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
			return
		}
		c.SetCookie("JWT", signedJWT, 3600*24, "/", "", false, true)

		util.WriteSuccessResponse(c, http.StatusOK, user, map[string]string{})
	}
}

func (S Server) UserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pass := c.Param("pass")
		if pass != "pass" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.Error{StatusCode: http.StatusUnauthorized, Message: "Unauthorized"})
			return
		}
		if c.Request.Method == http.MethodPost {
			var req model.User
			err := c.BindJSON(&req)
			if err != nil {
				util.WriteFailResponse(c, http.StatusBadRequest, util.Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
				return
			}

			createdUser, err := S.Manager.CreateUser(&req)

			util.WriteSuccessResponse(c, http.StatusOK, createdUser, nil)
		} else if c.Request.Method == http.MethodGet {
			users, _ := S.Manager.GetUsers()

			util.WriteSuccessResponse(c, http.StatusOK, users, nil)
		} else {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, util.Error{StatusCode: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		}
	}
}

func (S Server) GetTuitionHandler() middleware.MiddlewareHandlerFunc {
	return func(jwtS *util.JWTStruct) gin.HandlerFunc {
		return func(c *gin.Context) {
			tuition, err := S.Manager.GetTuitionByUsername(jwtS.Username)
			if err != nil {
				util.WriteFailResponse(c, http.StatusBadRequest, util.Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
			}

			util.WriteSuccessResponse(c, http.StatusOK, tuition, nil)
		}
	}
}
