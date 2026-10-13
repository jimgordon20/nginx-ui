package media

import (
	"github.com/0xJacky/Nginx-UI/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.RouterGroup) {
	r.GET("media/roots", ListRoots)
	r.GET("media/list", ListEntries)
	r.GET("media/file", GetFile)

	o := r.Group("", middleware.RequireSecureSession(), middleware.RejectInDemo())
	{
		o.POST("media/upload", Upload)
		o.POST("media/mkdir", Mkdir)
		o.POST("media/delete", Delete)
	}
}
