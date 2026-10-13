package nginx_log

import (
	"context"
	"net/http"
	"time"

	"github.com/0xJacky/Nginx-UI/internal/nginx_log"
	"github.com/0xJacky/Nginx-UI/internal/nginx_log/analytics"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
	"github.com/uozi-tech/cosy/logger"
)

// DomainAnalyticsRequest is the HTTP body for Domain Analytics.
type DomainAnalyticsRequest struct {
	LogPath   string `json:"log_path"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

// GetDomainAnalytics returns a composed Domain Analytics payload.
func GetDomainAnalytics(c *gin.Context) {
	var req DomainAnalyticsRequest
	if !cosy.BindAndValid(c, &req) {
		return
	}

	analyticsService := nginx_log.GetAnalytics()
	if analyticsService == nil {
		cosy.ErrHandler(c, nginx_log.ErrModernAnalyticsNotAvailable)
		return
	}

	safeLogPath := ""
	if req.LogPath != "" {
		decodedPath, err := decodeAndValidateLogPath(req.LogPath)
		if err != nil {
			cosy.ErrHandler(c, err)
			return
		}
		safeLogPath = decodedPath
	}

	if err := analyticsService.ValidateLogPath(safeLogPath); err != nil {
		cosy.ErrHandler(c, err)
		return
	}

	startTime := req.StartTime
	endTime := req.EndTime
	if startTime == 0 || endTime == 0 {
		now := time.Now().Unix()
		endTime = now
		startTime = now - 24*60*60
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	result, err := analyticsService.GetDomainAnalytics(ctx, &analytics.DomainAnalyticsRequest{
		LogPath:   safeLogPath,
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		logger.Debugf("Domain analytics failed: %v", err)
		cosy.ErrHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
