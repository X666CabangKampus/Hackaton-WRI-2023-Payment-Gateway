package signature

import (
	"github.com/gin-gonic/gin"
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

// WriteFailResponse creates error response for the http handler
func WriteFailResponse(c *gin.Context, statusCode int, error interface{}) {
	c.JSON(statusCode, error)
}
