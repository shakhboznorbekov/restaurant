package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func MakeRespond(c *gin.Context, data interface{}, err error) {
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
		"data":    data,
	})
}
