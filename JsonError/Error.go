package Error

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, err error) bool {

	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"path":      c.Request.URL.Path,
			"timestamp": time.Now(),
			"status":    http.StatusInternalServerError,
			"error":     http.StatusText(http.StatusInternalServerError),
			"message":   err.Error(),
		})
		return true
	}

	return false
}
