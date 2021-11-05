package Error

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, err error) bool {
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": false, "message": err.Error()})
		return true
	}
	return false
}
