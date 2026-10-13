package shield

import (
	"sort"

	"github.com/0xJacky/Nginx-UI/internal/sitecheck"
)

func OriginHealth() []HealthItem {
	service := sitecheck.GetService()
	if service == nil {
		return []HealthItem{}
	}
	sites := service.GetSites()
	items := make([]HealthItem, 0, len(sites))
	for _, info := range sites {
		if info == nil {
			continue
		}
		items = append(items, HealthItem{
			Name:       info.Name,
			URL:        info.GetURL(),
			Status:     info.Status,
			StatusCode: info.StatusCode,
			LatencyMS:  info.ResponseTime,
			Error:      info.Error,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	return items
}
