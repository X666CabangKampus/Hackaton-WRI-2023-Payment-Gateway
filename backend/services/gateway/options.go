package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/imrenagi/go-payment"
	"strconv"
)

// NewPaymentMethodListOptions accepts http.Request and returns set of option containing the price and its currency.
// Overwrite the
func NewPaymentMethodListOptions(c *gin.Context) ([]payment.Option, error) {
	var options []payment.Option
	var price float64
	var currency string
	var err error
	if len(c.PostForm("price")) > 0 {
		price, err = strconv.ParseFloat(c.PostForm("price"), 64)
		if err != nil {
			return nil, err
		}
	}
	if len(c.PostForm("currency")) > 0 {
		currency = c.PostForm("currency")
	}
	if price > 0 && currency != "" {
		options = append(options, payment.WithPrice(price, currency))
	}

	return options, nil
}
