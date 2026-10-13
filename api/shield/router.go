package shield

import (
	"github.com/0xJacky/Nginx-UI/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.RouterGroup) {
	r.GET("shield/health", OriginHealth)
	r.GET("shield/alerts", ListAlerts)
	r.GET("shield/:site", GetConfig)
	r.GET("shield/:site/preview", PreviewConfig)

	o := r.Group("", middleware.RequireSecureSession(), middleware.RejectInDemo())
	{
		o.POST("shield/:site", SaveConfig)
		o.POST("shield/:site/apply", ApplyConfig)
		o.POST("shield/:site/detach", DetachConfig)
		o.POST("shield/:site/preset/image", ApplyImagePreset)
		o.POST("shield/alerts", SaveAlert)
		o.POST("shield/alerts/delete", DeleteAlert)
		o.POST("shield/alerts/evaluate", EvaluateAlerts)
	}
}
