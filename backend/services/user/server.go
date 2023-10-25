package user

import (
	"backend-hacktober/services/gateway"
	"backend-hacktober/services/middleware"
	util "backend-hacktober/util"
	"github.com/gin-gonic/gin"
	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/invoice"
	"github.com/imrenagi/go-payment/manage"
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
			util.WriteFailResponse(c, http.StatusBadRequest, util.Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
			return
		}

		var user User
		s.DB.Where(map[string]any{
			"username": req.Username,
			"password": util.HashPassword(req.Password),
		}).Find(&user)

		if user.ID == 0 {
			util.WriteFailResponse(c, http.StatusBadRequest, util.Error{StatusCode: http.StatusBadRequest, Message: "Username or password is wrong"})
			return
		}

		signedJWT, err := util.CreateJWTSign(&util.JWTStruct{Username: user.Username})
		if err != nil {
			util.WriteFailResponse(c, http.StatusBadRequest, util.Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
			return
		}
		c.SetCookie("JWT", signedJWT, 3600*24, "", "", true, true)

		util.WriteSuccessResponse(c, http.StatusOK, user, map[string]string{})
	}
}

func (s Server) UserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pass := c.Param("pass")
		if pass != "pass" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.Error{StatusCode: http.StatusUnauthorized, Message: "Unauthorized"})
			return
		}
		if c.Request.Method == http.MethodPost {
			var req User
			err := c.BindJSON(&req)
			if err != nil {
				util.WriteFailResponse(c, http.StatusBadRequest, util.Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
				return
			}

			req.Password = util.HashPassword(req.Password)
			err = s.DB.Create(&req).Error
			if err != nil {
				util.WriteFailResponse(c, http.StatusBadRequest, util.Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
				return
			}

			util.WriteSuccessResponse(c, http.StatusOK, req, nil)
		} else if c.Request.Method == http.MethodGet {
			var users []User
			s.DB.Preload("Tuition").Find(&users)

			util.WriteSuccessResponse(c, http.StatusOK, users, nil)
		} else {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, util.Error{StatusCode: http.StatusMethodNotAllowed, Message: "Method not allowed"})
		}
	}
}

type PayTuitionsRequest struct {
	PaymentType      payment.PaymentType       `json:"payment_type"`
	Qty              int                       `json:"qty"`
	CreditCardDetail *invoice.CreditCardDetail `json:"credit_card,omitempty"`
}

func (s Server) PayTuitionHandler(paymentSrv *gateway.Server) middleware.MiddlewareHandlerFunc {
	return func(jwtS *util.JWTStruct) gin.HandlerFunc {
		return func(c *gin.Context) {
			var req PayTuitionsRequest
			err := c.BindJSON(&req)
			if err != nil {
				util.WriteFailResponse(c, http.StatusBadRequest, util.Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
				return
			}

			var user User
			s.DB.Find(&user, "username = ?", jwtS.Username)
			invoiceReq := manage.GenerateInvoiceRequest{
				Payment: struct {
					PaymentType      payment.PaymentType       `json:"payment_type"`
					CreditCardDetail *invoice.CreditCardDetail `json:"credit_card,omitempty"`
				}{
					PaymentType:      req.PaymentType,
					CreditCardDetail: req.CreditCardDetail,
				},
				Customer: struct {
					Name        string `json:"name"`
					Email       string `json:"email"`
					PhoneNumber string `json:"phone_number"`
				}{
					Name:        user.FullName,
					Email:       user.Email,
					PhoneNumber: user.Phone,
				},
				Items: []struct {
					Name         string  `json:"name"`
					Category     string  `json:"category"`
					MerchantName string  `json:"merchant"`
					Description  string  `json:"description"`
					Qty          int     `json:"qty"`
					Price        float64 `json:"price"`
					Currency     string  `json:"currency"`
				}{
					{
						Name:     "UKT",
						Category: "Tuition",
						Qty:      req.Qty,
						Price:    float64(user.TuitionFeeBase),
					},
				},
			}

			inv, err := paymentSrv.CreateInvoice(c.Copy(), &invoiceReq)
			if err != nil {
				gateway.WriteFailResponseFromError(c, err)
				return
			}

			numberOfSemester := len(user.Tuition)
			for i := 0; i < req.Qty; i++ {
				numberOfSemester++
				s.DB.Create(&UserTuitionFee{
					UserId:        user.ID,
					SemesterPay:   Semester(numberOfSemester),
					InvoiceNumber: inv.Number,
				})
			}

			util.WriteSuccessResponse(c, http.StatusOK, inv, nil)
		}
	}
}
