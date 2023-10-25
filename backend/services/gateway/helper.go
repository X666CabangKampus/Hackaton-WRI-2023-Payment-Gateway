package gateway

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/imrenagi/go-payment"
	"net/http"
)

type Error struct {
	StatusCode int    `json:"error_code"`
	Message    string `json:"error_message"`
}

type Meta struct {
	TotalItems  int    `json:"total_items"`
	TotalPages  int    `json:"total_pages"`
	CurrentPage int    `json:"cur_page"`
	Cursor      string `json:"last_cursor"`
}

// Empty used to return nothing
type Empty struct{}

func WriteSuccessResponse(c *gin.Context, statusCode int, data interface{}, headMap map[string]string) {
	if headMap != nil && len(headMap) > 0 {
		for key, val := range headMap {
			c.Header(key, val)
		}
	}
	c.JSON(statusCode, data)
}

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

	errorMsg := Error{
		Message:    err.Error(),
		StatusCode: statusCode,
	}

	c.JSON(statusCode, errorMsg)
}

// WriteFailResponse creates error response for the http handler
func WriteFailResponse(c *gin.Context, statusCode int, error interface{}) {
	c.JSON(statusCode, error)
}
