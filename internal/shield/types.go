package shield

type HotlinkConfig struct {
	Enabled        bool     `json:"enabled"`
	AllowEmpty     bool     `json:"allow_empty"`
	AllowReferers  []string `json:"allow_referers"`
	Extensions     []string `json:"extensions"`
}

type RateLimitConfig struct {
	Enabled   bool   `json:"enabled"`
	Rate      string `json:"rate"`
	Burst     int    `json:"burst"`
	PathRegex string `json:"path_regex"`
}

type WAFConfig struct {
	Enabled          bool     `json:"enabled"`
	BlockEmptyUA     bool     `json:"block_empty_ua"`
	BlockBadBots     bool     `json:"block_bad_bots"`
	BlockedUserAgents []string `json:"blocked_user_agents"`
	BlockedPaths     []string `json:"blocked_paths"`
}

type UnderAttackConfig struct {
	Enabled bool `json:"enabled"`
	Rate    string `json:"rate"`
	Burst   int    `json:"burst"`
}

type IPBlockConfig struct {
	Enabled bool     `json:"enabled"`
	DenyIPs []string `json:"deny_ips"`
}

type CacheConfig struct {
	Enabled    bool   `json:"enabled"`
	PathRegex  string `json:"path_regex"`
	Expires    string `json:"expires"`
	CacheControl string `json:"cache_control"`
}

type SecurityHeadersConfig struct {
	Enabled bool `json:"enabled"`
}

type Config struct {
	Site            string                `json:"site"`
	Applied         bool                  `json:"applied"`
	Hotlink         HotlinkConfig         `json:"hotlink"`
	RateLimit       RateLimitConfig       `json:"rate_limit"`
	WAF             WAFConfig             `json:"waf"`
	UnderAttack     UnderAttackConfig     `json:"under_attack"`
	IPBlock         IPBlockConfig         `json:"ip_block"`
	SecurityHeaders SecurityHeadersConfig `json:"security_headers"`
	Cache           CacheConfig           `json:"cache"`
	Preview         string                `json:"preview,omitempty"`
}

type AlertRule struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Enabled         bool   `json:"enabled"`
	Site            string `json:"site"`
	LogPath         string `json:"log_path"`
	Metric          string `json:"metric"`
	Threshold       float64 `json:"threshold"`
	WindowMinutes   int    `json:"window_minutes"`
	NotifyOnTrigger bool   `json:"notify_on_trigger"`
}

type AlertEvaluation struct {
	Rule    AlertRule `json:"rule"`
	Value   float64   `json:"value"`
	Triggered bool    `json:"triggered"`
	Message string    `json:"message"`
}

type HealthItem struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	LatencyMS  int64  `json:"latency_ms"`
	Error      string `json:"error,omitempty"`
}

func DefaultConfig(site string) Config {
	return Config{
		Site: site,
		Hotlink: HotlinkConfig{
			Enabled:    false,
			AllowEmpty: false,
			Extensions: []string{"jpg", "jpeg", "png", "gif", "webp", "avif", "css", "js"},
		},
		RateLimit: RateLimitConfig{
			Enabled:   false,
			Rate:      "10r/s",
			Burst:     20,
			PathRegex: `\.(jpg|jpeg|png|gif|webp|avif)$`,
		},
		WAF: WAFConfig{
			Enabled:      false,
			BlockEmptyUA: true,
			BlockBadBots: true,
			BlockedUserAgents: []string{
				"scrapy", "curl", "wget", "python-requests", "httpclient", "libwww",
			},
			BlockedPaths: []string{
				"/.env", "/wp-admin", "/xmlrpc.php", "/.git",
			},
		},
		UnderAttack: UnderAttackConfig{
			Enabled: false,
			Rate:    "2r/s",
			Burst:   5,
		},
		IPBlock: IPBlockConfig{
			Enabled: false,
			DenyIPs: []string{},
		},
		SecurityHeaders: SecurityHeadersConfig{Enabled: true},
		Cache: CacheConfig{
			Enabled:      false,
			PathRegex:    `\.(jpg|jpeg|png|gif|webp|avif|css|js|woff2?)$`,
			Expires:      "30d",
			CacheControl: "public, max-age=2592000, immutable",
		},
	}
}

func ImageSitePreset(site string) Config {
	cfg := DefaultConfig(site)
	cfg.Hotlink.Enabled = true
	cfg.RateLimit.Enabled = true
	cfg.WAF.Enabled = true
	cfg.SecurityHeaders.Enabled = true
	cfg.Cache.Enabled = true
	return cfg
}
