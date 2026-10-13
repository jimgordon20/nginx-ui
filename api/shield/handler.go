package shield

import (
	"net/http"

	"github.com/0xJacky/Nginx-UI/internal/shield"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
)

func GetConfig(c *gin.Context) {
	site := c.Param("site")
	cfg, err := shield.Load(site)
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, cfg)
}

func PreviewConfig(c *gin.Context) {
	site := c.Param("site")
	cfg, err := shield.Load(site)
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"preview": shield.Render(cfg),
		"applied": cfg.Applied,
	})
}

func SaveConfig(c *gin.Context) {
	site := c.Param("site")
	var cfg shield.Config
	if !cosy.BindAndValid(c, &cfg) {
		return
	}
	cfg.Site = site
	if err := shield.Save(cfg); err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	if err := shield.WriteRendered(cfg); err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	cfg.Preview = shield.Render(cfg)
	cfg.Applied = shield.IsApplied(site)
	c.JSON(http.StatusOK, cfg)
}

func ApplyConfig(c *gin.Context) {
	site := c.Param("site")
	var cfg shield.Config
	if !cosy.BindAndValid(c, &cfg) {
		return
	}
	cfg.Site = site
	if err := shield.Apply(c.Request.Context(), cfg); err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	loaded, _ := shield.Load(site)
	c.JSON(http.StatusOK, loaded)
}

func DetachConfig(c *gin.Context) {
	site := c.Param("site")
	if err := shield.Detach(c.Request.Context(), site); err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	loaded, _ := shield.Load(site)
	c.JSON(http.StatusOK, loaded)
}

func ApplyImagePreset(c *gin.Context) {
	site := c.Param("site")
	cfg := shield.ImageSitePreset(site)
	if err := shield.Apply(c.Request.Context(), cfg); err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	loaded, _ := shield.Load(site)
	c.JSON(http.StatusOK, loaded)
}

func OriginHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": shield.OriginHealth()})
}

func ListAlerts(c *gin.Context) {
	rules, err := shield.LoadAlerts()
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rules})
}

func SaveAlert(c *gin.Context) {
	var rule shield.AlertRule
	if !cosy.BindAndValid(c, &rule) {
		return
	}
	rules, err := shield.UpsertAlert(rule)
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rules})
}

func DeleteAlert(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
	}
	if !cosy.BindAndValid(c, &req) {
		return
	}
	rules, err := shield.DeleteAlert(req.ID)
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rules})
}

func EvaluateAlerts(c *gin.Context) {
	results, err := shield.EvaluateAlerts(c.Request.Context())
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results})
}
