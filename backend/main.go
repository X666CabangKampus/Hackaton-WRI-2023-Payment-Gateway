package main

import (
	"backend-hacktober/services"
	"backend-hacktober/services/email"
	"github.com/gin-gonic/gin"
	"github.com/imrenagi/go-payment/util/localconfig"
	"os"
	"strconv"
)

func getSecret() (*localconfig.Secret, error) {
	secret, err := localconfig.LoadSecret("./conf/secret.yaml")

	if secret.DB.Host == "" || secret.DB.UserName == "" || secret.DB.Password == "" || secret.DB.DBName == "" || secret.DB.Port == 0 {
		secret.DB.Host = os.Getenv("DB_HOST")
		secret.DB.UserName = os.Getenv("DB_USERNAME")
		secret.DB.Password = os.Getenv("DB_PASSWORD")
		secret.DB.DBName = os.Getenv("DB_NAME")

		secret.DB.Port, err = strconv.Atoi(os.Getenv("DB_PORT"))
		if err != nil {
			return nil, err
		}
	}

	if secret.Payment.Midtrans.SecretKey == "" || secret.Payment.Midtrans.ClientKey == "" {
		secret.Payment.Midtrans.SecretKey = os.Getenv("MIDTRANS_SECRET_KEY")
		secret.Payment.Midtrans.ClientKey = os.Getenv("MIDTRANS_CLIENT_KEY")
	}

	return secret, nil
}

func main() {
	r := gin.Default()

	config, err := localconfig.LoadConfig("./conf/config.yaml")
	if err != nil {
		panic(err)
	}

	secret, err := getSecret()
	if err != nil {
		panic(err)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello world!"})
	})
	services.NewSrv(r, config, secret).Routes()

	go emailservice.EmailQueueReceiver()
	err = r.Run("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
}
