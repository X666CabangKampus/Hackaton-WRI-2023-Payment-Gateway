package services

import (
	"backend-hacktober/modules"
	srvPayment "backend-hacktober/services/gateway"
	"backend-hacktober/services/middleware"
	srvUser "backend-hacktober/services/user"
	userManage "backend-hacktober/services/user/manage"
	"backend-hacktober/services/user/model"
	"backend-hacktober/services/user/repository"
	util "backend-hacktober/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/imrenagi/go-payment/datastore/inmemory"
	dssql "github.com/imrenagi/go-payment/datastore/sql"
	"github.com/imrenagi/go-payment/gateway/midtrans"
	"github.com/imrenagi/go-payment/invoice"
	paymentManage "github.com/imrenagi/go-payment/manage"
	"github.com/imrenagi/go-payment/subscription"
	"github.com/imrenagi/go-payment/util/localconfig"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

type Srv struct {
	Router     *gin.Engine
	PaymentSrv *srvPayment.Server
	UserSrv    *srvUser.Server
}

func NewSrv(router *gin.Engine, config *localconfig.Config, secret *localconfig.Secret) *Srv {

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta", secret.DB.Host, secret.DB.UserName, secret.DB.Password, secret.DB.DBName, secret.DB.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	err = util.AutoOrderedMigrate(
		db,
		subscription.Schedule{},
		subscription.Subscription{},
		invoice.Invoice{},
		invoice.LineItem{},
		midtrans.TransactionStatus{},
		invoice.Payment{},
		invoice.BillingAddress{},
		invoice.CreditCardDetail{},
		model.User{},
		model.UserTuitionFee{},
	)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	mPayment := paymentManage.NewManager(*config, secret.Payment)
	mPayment.MustMidtransTransactionStatusRepository(dssql.NewMidtransTransactionRepository(db))
	mPayment.MustInvoiceRepository(dssql.NewInvoiceRepository(db))
	mPayment.MustSubscriptionRepository(dssql.NewSubscriptionRepository(db))
	mPayment.MustPaymentConfigReader(inmemory.NewPaymentConfigRepository("conf/payment-methods.yaml"))

	mUser := userManage.NewManager(repository.NewUserRepository(db))

	return &Srv{
		Router:     router,
		PaymentSrv: srvPayment.NewServer(mPayment),
		UserSrv:    srvUser.NewServer(*mUser),
	}
}

func (S Srv) PayTuitionHandler() middleware.MiddlewareHandlerFunc {
	return func(jwtS *util.JWTStruct) gin.HandlerFunc {
		return func(c *gin.Context) {
			var req userManage.PayTuitionsRequest
			err := c.BindJSON(&req)
			if err != nil {
				util.WriteFailResponse(c, http.StatusBadRequest, util.Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
				return
			}

			inv, err := S.UserSrv.Manager.PayTuition(S.PaymentSrv, c.Request.Context(), jwtS.Username, &req)
			if err != nil {
				util.WriteFailResponse(c, http.StatusBadRequest, util.Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
				return
			}

			util.WriteSuccessResponse(c, http.StatusOK, inv, nil)
		}
	}
}

func (S Srv) MidtransTransactionCallbackWrapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		var notification coreapi.TransactionStatusResponse
		err := c.BindJSON(&notification)
		if err != nil {
			util.WriteFailResponse(c, http.StatusBadRequest, util.Error{
				StatusCode: http.StatusBadRequest,
				Message:    "Request can't be parsed",
			})
			return
		}
		S.PaymentSrv.MidtransTransactionCallbackHandler(&notification)(c)

		user, err := S.UserSrv.Manager.GetUserFromInvoiceNumber(notification.OrderID)
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}

		modules.SendActivationMail(user.Email, user.FullName, "Your payment has been processed")

		util.WriteSuccessResponse(c, http.StatusOK, nil, nil)
	}

}

func (S Srv) Routes() {
	S.Router.GET("/payment/methods", S.PaymentSrv.GetPaymentMethodsHandler())
	S.Router.POST("/payment/invoices", middleware.MiddlewareJWT(S.PaymentSrv.CreateInvoiceHandler()))
	S.Router.POST("/payment/midtrans/callback", S.MidtransTransactionCallbackWrapper())
	S.Router.GET("/payment/midtrans/callback", func(context *gin.Context) {
		util.WriteSuccessResponse(context, http.StatusOK, nil, nil)
	})
	S.Router.POST("/login", S.UserSrv.LoginHandler())
	S.Router.POST("/tuition/pay", middleware.MiddlewareJWT(S.PayTuitionHandler()))
	S.Router.GET("/tuition", middleware.MiddlewareJWT(S.UserSrv.GetTuitionHandler()))
	S.Router.POST("/user/:pass", S.UserSrv.UserHandler())
	S.Router.GET("/user/:pass", S.UserSrv.UserHandler())
}
