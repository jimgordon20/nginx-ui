package analytics

import (
	"path"
	"sort"
	"strings"
)

// ExtensionAccessStats represents request volume grouped by URL file extension.
type ExtensionAccessStats struct {
	Extension string  `json:"extension"`
	Count     int     `json:"count"`
	Percent   float64 `json:"percent"`
}

// AggregateExtensions groups path_exact facet terms by file extension.
// Paths without a recognizable extension are bucketed as "(none)".
func AggregateExtensions(paths []KeyValue, limit int) []ExtensionAccessStats {
	if limit <= 0 {
		limit = 10
	}

	totals := make(map[string]int)
	var totalCount int
	for _, item := range paths {
		ext := fileExtension(item.Key)
		totals[ext] += item.Value
		totalCount += item.Value
	}

	result := make([]ExtensionAccessStats, 0, len(totals))
	for ext, count := range totals {
		percent := 0.0
		if totalCount > 0 {
			percent = float64(count) / float64(totalCount) * 100
		}
		result = append(result, ExtensionAccessStats{
			Extension: ext,
			Count:     count,
			Percent:   percent,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Extension < result[j].Extension
		}
		return result[i].Count > result[j].Count
	})

	if len(result) > limit {
		result = result[:limit]
	}
	return result
}

func fileExtension(rawPath string) string {
	cleaned := strings.TrimSpace(rawPath)
	if cleaned == "" {
		return "(none)"
	}

	// Strip query string / fragment if present in indexed path values.
	if idx := strings.IndexAny(cleaned, "?#"); idx >= 0 {
		cleaned = cleaned[:idx]
	}

	base := path.Base(cleaned)
	if base == "." || base == "/" || base == "" {
		return "(none)"
	}

	ext := path.Ext(base)
	if ext == "" || ext == "." {
		return "(none)"
	}

	return strings.ToLower(ext)
}
