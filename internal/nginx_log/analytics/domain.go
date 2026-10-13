package analytics

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/0xJacky/Nginx-UI/internal/nginx_log/searcher"
)

// DomainAnalyticsRequest is the query for the Domain Analytics dashboard.
type DomainAnalyticsRequest struct {
	LogPath   string
	StartTime int64
	EndTime   int64
}

// DomainLogEntry is a compact recent-request row for the Domain Analytics table.
type DomainLogEntry struct {
	Timestamp int64  `json:"timestamp"`
	IP        string `json:"ip"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	BytesSent int64  `json:"bytes_sent"`
	Raw       string `json:"raw"`
}

// DomainBlockedIPStat ranks client IPs that received 401/403 responses.
type DomainBlockedIPStat struct {
	IP        string `json:"ip"`
	Count     int    `json:"count"`
	Status401 int    `json:"status_401"`
	Status403 int    `json:"status_403"`
}

// DomainAnalytics is the composed payload for the Domain Analytics UI.
type DomainAnalytics struct {
	HourlyStats     []HourlyAccessStats    `json:"hourly_stats"`
	DailyStats      []DailyAccessStats     `json:"daily_stats"`
	TopURLs         []URLAccessStats       `json:"top_urls"`
	Summary         DashboardSummary       `json:"summary"`
	StatusCodeDist  map[string]int         `json:"status_code_distribution"`
	Extensions      []ExtensionAccessStats `json:"extensions"`
	SuccessRate     float64                `json:"success_rate"`
	ClientErrorRate float64                `json:"client_error_rate"`
	ServerErrorRate float64                `json:"server_error_rate"`
	TotalStatusHits int                    `json:"total_status_hits"`
	RecentEntries   []DomainLogEntry       `json:"recent_entries"`
	SecurityEvents  []DomainLogEntry       `json:"security_events"`
	TopBlockedIPs   []DomainBlockedIPStat  `json:"top_blocked_ips"`
	BotRequests     int                    `json:"bot_requests"`
	HumanRequests   int                    `json:"human_requests"`
	BotPercent      float64                `json:"bot_percent"`
}

// GetDomainAnalytics composes dashboard series with status and extension breakdowns.
// An empty LogPath means all indexed access logs (All Domains).
func (s *service) GetDomainAnalytics(ctx context.Context, req *DomainAnalyticsRequest) (*DomainAnalytics, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if err := s.ValidateTimeRange(req.StartTime, req.EndTime); err != nil {
		return nil, fmt.Errorf("invalid time range: %w", err)
	}
	if err := s.ValidateLogPath(req.LogPath); err != nil {
		return nil, err
	}

	dashboardReq := &DashboardQueryRequest{
		LogPath:   req.LogPath,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	if req.LogPath != "" {
		dashboardReq.LogPaths = []string{req.LogPath}
	}

	dashboard, err := s.GetDashboardAnalytics(ctx, dashboardReq)
	if err != nil {
		return nil, err
	}

	facetReq := &searcher.SearchRequest{
		StartTime:      &req.StartTime,
		EndTime:        &req.EndTime,
		UseCache:       true,
		IncludeFacets:  true,
		IncludeStats:   true,
		FacetFields:    []string{"status", "path_exact"},
		FacetSize:      200,
		UseMainLogPath: true,
		SortBy:         "timestamp",
		SortOrder:      "desc",
		Limit:          -1,
	}
	if req.LogPath != "" {
		facetReq.LogPaths = []string{req.LogPath}
	}

	entries, err := s.GetLogEntriesStats(ctx, facetReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get status and path facets: %w", err)
	}

	statusDist := entries.StatusCodeDist
	if statusDist == nil {
		statusDist = make(map[string]int)
	}

	successRate, clientErrorRate, serverErrorRate, totalStatusHits := summarizeStatusRates(statusDist)

	recent, err := s.fetchRecentDomainEntries(ctx, req)
	if err != nil {
		return nil, err
	}

	securityEvents, topBlockedIPs, err := s.fetchSecurityBouncer(ctx, req)
	if err != nil {
		return nil, err
	}

	botRequests, humanRequests, botPercent := botHumanSplit(dashboard.Browsers, dashboard.Summary.TotalPV)

	return &DomainAnalytics{
		HourlyStats:     dashboard.HourlyStats,
		DailyStats:      dashboard.DailyStats,
		TopURLs:         dashboard.TopURLs,
		Summary:         dashboard.Summary,
		StatusCodeDist:  statusDist,
		Extensions:      AggregateExtensions(entries.TopPaths, 12),
		SuccessRate:     successRate,
		ClientErrorRate: clientErrorRate,
		ServerErrorRate: serverErrorRate,
		TotalStatusHits: totalStatusHits,
		RecentEntries:   recent,
		SecurityEvents:  securityEvents,
		TopBlockedIPs:   topBlockedIPs,
		BotRequests:     botRequests,
		HumanRequests:   humanRequests,
		BotPercent:      botPercent,
	}, nil
}

func botHumanSplit(browsers []BrowserAccessStats, totalPV int) (bot, human int, percent float64) {
	for _, item := range browsers {
		if strings.EqualFold(item.Browser, "Bot") || strings.Contains(strings.ToLower(item.Browser), "bot") {
			bot += item.Count
			continue
		}
		human += item.Count
	}
	if totalPV > 0 && bot+human == 0 {
		human = totalPV
	}
	counted := bot + human
	if counted > 0 {
		percent = float64(bot) / float64(counted) * 100
	}
	return bot, human, percent
}

func (s *service) fetchRecentDomainEntries(ctx context.Context, req *DomainAnalyticsRequest) ([]DomainLogEntry, error) {
	return s.fetchDomainEntries(ctx, req, nil, 40)
}

func (s *service) fetchSecurityBouncer(ctx context.Context, req *DomainAnalyticsRequest) ([]DomainLogEntry, []DomainBlockedIPStat, error) {
	// Pull a wider sample so Top 5 blocked IPs is meaningful, then keep a
	// shorter table for the Security & Bouncer panel.
	hits, err := s.fetchDomainEntries(ctx, req, []int{401, 403}, 250)
	if err != nil {
		return nil, nil, err
	}

	events := hits
	if len(events) > 40 {
		events = events[:40]
	}

	type agg struct {
		count     int
		status401 int
		status403 int
	}
	byIP := make(map[string]*agg)
	for _, entry := range hits {
		ip := entry.IP
		if ip == "" {
			ip = "(unknown)"
		}
		item, ok := byIP[ip]
		if !ok {
			item = &agg{}
			byIP[ip] = item
		}
		item.count++
		switch entry.Status {
		case 401:
			item.status401++
		case 403:
			item.status403++
		}
	}

	ranked := make([]DomainBlockedIPStat, 0, len(byIP))
	for ip, item := range byIP {
		ranked = append(ranked, DomainBlockedIPStat{
			IP:        ip,
			Count:     item.count,
			Status401: item.status401,
			Status403: item.status403,
		})
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Count == ranked[j].Count {
			return ranked[i].IP < ranked[j].IP
		}
		return ranked[i].Count > ranked[j].Count
	})
	if len(ranked) > 5 {
		ranked = ranked[:5]
	}

	return events, ranked, nil
}

func (s *service) fetchDomainEntries(ctx context.Context, req *DomainAnalyticsRequest, statusCodes []int, limit int) ([]DomainLogEntry, error) {
	if limit <= 0 {
		limit = 40
	}
	searchReq := &searcher.SearchRequest{
		StartTime:      &req.StartTime,
		EndTime:        &req.EndTime,
		UseCache:       true,
		UseMainLogPath: true,
		SortBy:         "timestamp",
		SortOrder:      "desc",
		Limit:          limit,
		Offset:         0,
		StatusCodes:    statusCodes,
	}
	if req.LogPath != "" {
		searchReq.LogPaths = []string{req.LogPath}
	}

	result, err := s.searcher.Search(ctx, searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch domain entries: %w", err)
	}

	entries := make([]DomainLogEntry, 0, len(result.Hits))
	for _, hit := range result.Hits {
		if hit == nil || hit.Fields == nil {
			continue
		}
		entries = append(entries, DomainLogEntry{
			Timestamp: int64Field(hit.Fields, "timestamp"),
			IP:        stringField(hit.Fields, "ip"),
			Method:    stringField(hit.Fields, "method"),
			Path:      stringField(hit.Fields, "path"),
			Status:    intField(hit.Fields, "status"),
			BytesSent: int64Field(hit.Fields, "bytes_sent"),
			Raw:       stringField(hit.Fields, "raw"),
		})
	}
	return entries, nil
}

func stringField(fields map[string]interface{}, key string) string {
	if v, ok := fields[key].(string); ok {
		return v
	}
	return ""
}

func intField(fields map[string]interface{}, key string) int {
	switch v := fields[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	default:
		return 0
	}
}

func int64Field(fields map[string]interface{}, key string) int64 {
	switch v := fields[key].(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	default:
		return 0
	}
}

func summarizeStatusRates(dist map[string]int) (success, clientErr, serverErr float64, total int) {
	var ok2xx, err4xx, err5xx int
	for code, count := range dist {
		total += count
		switch {
		case len(code) > 0 && code[0] == '2':
			ok2xx += count
		case len(code) > 0 && code[0] == '4':
			err4xx += count
		case len(code) > 0 && code[0] == '5':
			err5xx += count
		}
	}
	if total == 0 {
		return 0, 0, 0, 0
	}
	return float64(ok2xx) / float64(total) * 100,
		float64(err4xx) / float64(total) * 100,
		float64(err5xx) / float64(total) * 100,
		total
}
