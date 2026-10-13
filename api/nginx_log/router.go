package nginx_log

import (
	"github.com/0xJacky/Nginx-UI/internal/middleware"
	"github.com/gin-gonic/gin"
)

// InitRouter registers all the nginx log related routes
func InitRouter(r *gin.RouterGroup) {

	r.GET("nginx_logs", GetLogList)
	r.POST("nginx_log/page", GetNginxLogPage)
	r.POST("nginx_log/analytics", GetLogAnalytics)
	r.GET("nginx_log/entries", GetLogEntries)
	r.POST("nginx_log/search", AdvancedSearchLogs)
	r.GET("nginx_log/preflight", GetLogPreflight)
	r.POST("nginx_log/dashboard", GetDashboardAnalytics)
	r.POST("nginx_log/domain_analytics", GetDomainAnalytics)
	r.POST("nginx_log/geo/world", GetWorldMapData)
	r.POST("nginx_log/geo/china", GetChinaMapData)
	r.POST("nginx_log/geo/china/city", GetChinaCityMapData)
	r.POST("nginx_log/geo/stats", GetGeoStats)
	r.POST("nginx_log/index/rebuild", RebuildIndex)
	r.POST("nginx_log/settings/advanced_indexing/enable", EnableAdvancedIndexing)
	r.POST("nginx_log/settings/advanced_indexing/disable", DisableAdvancedIndexing)
	r.GET("nginx_log/settings/advanced_indexing/status", GetAdvancedIndexingStatus)
	r.GET("nginx_log/default_log_dir", GetDefaultLogDir)

	// Destructive: truncate a specific domain log file in place.
	o := r.Group("", middleware.RequireSecureSession(), middleware.RejectInDemo())
	{
		o.POST("nginx_log/clear", ClearDomainLog)
	}
}

func InitPublicRouter(r *gin.RouterGroup) {
	r.GET("nginx_log/geo/boundary/:filename", GetGeoBoundaryFile)
}

func InitWebSocketRouter(r *gin.RouterGroup) {
	r.GET("nginx_log", Log)
}
