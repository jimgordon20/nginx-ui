package nginx_log

import (
	"net/http"

	"github.com/0xJacky/Nginx-UI/internal/nginx_log"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
)

// ClearDomainLogRequest is the body for truncating a domain log file.
type ClearDomainLogRequest struct {
	LogPath string `json:"log_path" binding:"required"`
}

// ClearDomainLog truncates the selected domain access log in place.
func ClearDomainLog(c *gin.Context) {
	var req ClearDomainLogRequest
	if !cosy.BindAndValid(c, &req) {
		return
	}

	safeLogPath, err := decodeAndValidateLogPath(req.LogPath)
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	if safeLogPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "log_path is required",
		})
		return
	}

	if err := nginx_log.ClearDomainLog(safeLogPath); err != nil {
		cosy.ErrHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "log cleared",
		"log_path": safeLogPath,
	})
}
