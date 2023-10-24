package gateway

import (
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
	PaymentSrv *Server
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
	db.AutoMigrate(
		&subscription.Schedule{},
		&subscription.Subscription{},
		&midtrans.TransactionStatus{},
		&invoice.CreditCardDetail{},
		&invoice.LineItem{},
		&invoice.Payment{},
		&invoice.BillingAddress{},
		&invoice.Invoice{},
	)

	m := manage.NewManager(*config, secret.Payment)
	m.MustMidtransTransactionStatusRepository(dssql.NewMidtransTransactionRepository(db))
	m.MustInvoiceRepository(dssql.NewInvoiceRepository(db))
	m.MustSubscriptionRepository(dssql.NewSubscriptionRepository(db))
	m.MustPaymentConfigReader(inmemory.NewPaymentConfigRepository("conf/payment-methods.yaml"))

	return &Srv{
		Router:     router,
		PaymentSrv: NewServer(m),
	}
}

func (S Srv) Routes() {
	S.Router.GET("/payment/methods", S.PaymentSrv.GetPaymentMethodsHandler())
	S.Router.POST("/payment/invoices", S.PaymentSrv.CreateInvoiceHandler())
	S.Router.POST("/payment/midtrans/callback", S.PaymentSrv.MidtransTransactionCallbackHandler())
	S.Router.POST("/payment/subscriptions", S.PaymentSrv.CreateSubscriptionHandler())
	S.Router.POST("/payment/subscriptions/:subscription_number/pause", S.PaymentSrv.PauseSubscriptionHandler())
	S.Router.POST("/payment/subscriptions/:subscription_number/stop", S.PaymentSrv.StopSubscriptionHandler())
	S.Router.POST("/payment/subscriptions/:subscription_number/resume", S.PaymentSrv.ResumeSubscriptionHandler())
	S.Router.PUT("/payment/subscriptions/:subscription_number/pause", S.PaymentSrv.PauseSubscriptionHandler())
	S.Router.PUT("/payment/subscriptions/:subscription_number/stop", S.PaymentSrv.StopSubscriptionHandler())
	S.Router.PUT("/payment/subscriptions/:subscription_number/resume", S.PaymentSrv.ResumeSubscriptionHandler())
}
