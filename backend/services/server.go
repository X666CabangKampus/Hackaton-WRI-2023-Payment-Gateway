package services

import (
	"backend-hacktober/services/gateway"
	"backend-hacktober/services/middleware"
	"backend-hacktober/services/user"
	util "backend-hacktober/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/imrenagi/go-payment/datastore/inmemory"
	dssql "github.com/imrenagi/go-payment/datastore/sql"
	"github.com/imrenagi/go-payment/gateway/midtrans"
	"github.com/imrenagi/go-payment/invoice"
	"github.com/imrenagi/go-payment/manage"
	"github.com/imrenagi/go-payment/subscription"
	"github.com/imrenagi/go-payment/util/localconfig"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Srv struct {
	Router     *gin.Engine
	PaymentSrv *gateway.Server
	UserSrv    *user.Server
}

func NewSrv(router *gin.Engine) *Srv {
	config, err := localconfig.LoadConfig("./conf/config.yaml")
	if err != nil {
		panic(err)
	}

	secret, err := localconfig.LoadSecret("./conf/secret.yaml")
	if err != nil {
		panic(err)
	}

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
		user.User{},
		user.UserTuitionFee{},
	)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	m := manage.NewManager(*config, secret.Payment)
	m.MustMidtransTransactionStatusRepository(dssql.NewMidtransTransactionRepository(db))
	m.MustInvoiceRepository(dssql.NewInvoiceRepository(db))
	m.MustSubscriptionRepository(dssql.NewSubscriptionRepository(db))
	m.MustPaymentConfigReader(inmemory.NewPaymentConfigRepository("conf/payment-methods.yaml"))

	return &Srv{
		Router:     router,
		PaymentSrv: gateway.NewServer(m),
		UserSrv:    user.NewServer(db),
	}
}

func (S Srv) Routes() {
	S.Router.GET("/payment/methods", S.PaymentSrv.GetPaymentMethodsHandler())
	S.Router.POST("/payment/invoices", middleware.MiddlewareJWT(S.PaymentSrv.CreateInvoiceHandler()))
	S.Router.POST("/payment/midtrans/callback", S.PaymentSrv.MidtransTransactionCallbackHandler())
	S.Router.POST("/login", S.UserSrv.LoginHandler())
	S.Router.POST("/pay-tuition", middleware.MiddlewareJWT(S.UserSrv.PayTuitionHandler(S.PaymentSrv)))
	S.Router.POST("/user/:pass", S.UserSrv.UserHandler())
	S.Router.GET("/user/:pass", S.UserSrv.UserHandler())
}
