package gateway

import (
	util "backend-hacktober/util"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/imrenagi/go-payment"
	"net/http"
)

func WriteFailResponseFromError(c *gin.Context, err error) {
	var statusCode int
	if errors.Is(err, payment.ErrNotFound) {
		statusCode = http.StatusNotFound
	} else if errors.Is(err, payment.ErrInternal) {
		statusCode = http.StatusInternalServerError
	} else if errors.Is(err, payment.ErrDatabase) {
		statusCode = http.StatusInternalServerError
	} else if errors.Is(err, payment.ErrBadRequest) {
		statusCode = http.StatusBadRequest
	} else if errors.Is(err, payment.ErrCantProceed) {
		statusCode = http.StatusUnprocessableEntity
	} else if errors.Is(err, payment.ErrUnauthorized) {
		statusCode = http.StatusUnauthorized
	} else if errors.Is(err, payment.ErrForbidden) {
		statusCode = http.StatusForbidden
	} else {
		statusCode = http.StatusInternalServerError
	}

	errorMsg := util.Error{
		Message:    err.Error(),
		StatusCode: statusCode,
	}

	c.JSON(statusCode, errorMsg)
}
