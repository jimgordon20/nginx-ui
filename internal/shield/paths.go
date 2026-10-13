package shield

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/0xJacky/Nginx-UI/internal/nginx"
)

const (
	IncludeDir = "nginx-ui/shield"
	FileExt    = ".conf"
	StateExt   = ".json"
	ZonesFile  = "zones.conf"
)

var siteNamePattern = regexp.MustCompile(`^[A-Za-z0-9._@-]+$`)

func Dir() string {
	return nginx.GetConfPath(IncludeDir)
}

func SiteConfPath(site string) string {
	return nginx.GetConfPath(IncludeDir, sanitizeSite(site)+FileExt)
}

func SiteStatePath(site string) string {
	return nginx.GetConfPath(IncludeDir, sanitizeSite(site)+StateExt)
}

func ZonesPath() string {
	return nginx.GetConfPath(IncludeDir, ZonesFile)
}

func IncludePath(site string) string {
	return IncludeDir + "/" + sanitizeSite(site) + FileExt
}

func ZonesIncludePath() string {
	return IncludeDir + "/" + ZonesFile
}

func AlertsPath() string {
	return nginx.GetConfPath(IncludeDir, "alerts.json")
}

func ValidSiteName(name string) bool {
	return name != "" && siteNamePattern.MatchString(name) && !strings.Contains(name, "..")
}

func sanitizeSite(site string) string {
	site = filepath.Base(strings.TrimSpace(site))
	site = strings.ReplaceAll(site, string(filepath.Separator), "_")
	return site
}
