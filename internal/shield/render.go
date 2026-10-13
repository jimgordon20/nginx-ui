package shield

import (
	"fmt"
	"strings"
)

func Render(cfg Config) string {
	var b strings.Builder
	b.WriteString("# Managed by NGINX UI Site Shield — do not edit by hand\n")

	if cfg.SecurityHeaders.Enabled {
		b.WriteString(`add_header X-Content-Type-Options "nosniff" always;
add_header X-Frame-Options "SAMEORIGIN" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
add_header X-XSS-Protection "1; mode=block" always;
`)
	}

	if cfg.IPBlock.Enabled {
		for _, ip := range cfg.IPBlock.DenyIPs {
			ip = strings.TrimSpace(ip)
			if ip == "" {
				continue
			}
			fmt.Fprintf(&b, "deny %s;\n", ip)
		}
	}

	if cfg.WAF.Enabled {
		if cfg.WAF.BlockEmptyUA {
			b.WriteString(`if ($http_user_agent = "") { return 403; }
`)
		}
		if cfg.WAF.BlockBadBots {
			for _, agent := range cfg.WAF.BlockedUserAgents {
				agent = strings.TrimSpace(agent)
				if agent == "" {
					continue
				}
				fmt.Fprintf(&b, "if ($http_user_agent ~* %s) { return 403; }\n", regexpQuote(agent))
			}
		}
		for _, path := range cfg.WAF.BlockedPaths {
			path = strings.TrimSpace(path)
			if path == "" {
				continue
			}
			fmt.Fprintf(&b, "location = %s { return 404; }\n", path)
		}
	}

	if cfg.UnderAttack.Enabled {
		zone := zoneName(cfg.Site, "attack")
		fmt.Fprintf(&b, "limit_req zone=%s burst=%d nodelay;\n", zone, nonzero(cfg.UnderAttack.Burst, 5))
	}

	if needsMediaLocation(cfg) {
		b.WriteString(renderMediaLocation(cfg))
	}

	return b.String()
}

func needsMediaLocation(cfg Config) bool {
	return cfg.Hotlink.Enabled || cfg.RateLimit.Enabled || cfg.Cache.Enabled
}

func mediaPathRegex(cfg Config) string {
	if cfg.Hotlink.Enabled {
		exts := cfg.Hotlink.Extensions
		if len(exts) == 0 {
			exts = []string{"jpg", "jpeg", "png", "gif", "webp", "avif"}
		}
		return `\.(` + strings.Join(exts, "|") + `)$`
	}
	if cfg.RateLimit.Enabled && strings.TrimSpace(cfg.RateLimit.PathRegex) != "" {
		return cfg.RateLimit.PathRegex
	}
	if cfg.Cache.Enabled && strings.TrimSpace(cfg.Cache.PathRegex) != "" {
		return cfg.Cache.PathRegex
	}
	return `\.(jpg|jpeg|png|gif|webp|avif|css|js)$`
}

func renderMediaLocation(cfg Config) string {
	var b strings.Builder
	fmt.Fprintf(&b, "location ~* %s {\n", mediaPathRegex(cfg))

	if cfg.RateLimit.Enabled {
		fmt.Fprintf(&b, "    limit_req zone=%s burst=%d nodelay;\n",
			zoneName(cfg.Site, "media"), nonzero(cfg.RateLimit.Burst, 20))
	}

	if cfg.Hotlink.Enabled {
		referers := "server_names"
		if cfg.Hotlink.AllowEmpty {
			referers = "none blocked server_names"
		}
		for _, ref := range cfg.Hotlink.AllowReferers {
			ref = strings.TrimSpace(ref)
			if ref != "" {
				referers += " " + ref
			}
		}
		fmt.Fprintf(&b, "    valid_referers %s;\n", referers)
		b.WriteString("    if ($invalid_referer) { return 403; }\n")
	}

	if cfg.Cache.Enabled {
		expires := cfg.Cache.Expires
		if expires == "" {
			expires = "30d"
		}
		cc := cfg.Cache.CacheControl
		if cc == "" {
			cc = "public, max-age=2592000"
		}
		fmt.Fprintf(&b, "    expires %s;\n", expires)
		fmt.Fprintf(&b, "    add_header Cache-Control \"%s\" always;\n", cc)
		b.WriteString("    access_log off;\n")
	}

	b.WriteString("}\n")
	return b.String()
}

func RenderZones(configs []Config) string {
	var b strings.Builder
	b.WriteString("# Managed by NGINX UI Site Shield — http-context rate limit zones\n")
	seen := map[string]bool{}
	for _, cfg := range configs {
		if cfg.RateLimit.Enabled {
			name := zoneName(cfg.Site, "media")
			if !seen[name] {
				rate := cfg.RateLimit.Rate
				if rate == "" {
					rate = "10r/s"
				}
				fmt.Fprintf(&b, "limit_req_zone $binary_remote_addr zone=%s:10m rate=%s;\n", name, rate)
				seen[name] = true
			}
		}
		if cfg.UnderAttack.Enabled {
			name := zoneName(cfg.Site, "attack")
			if !seen[name] {
				rate := cfg.UnderAttack.Rate
				if rate == "" {
					rate = "2r/s"
				}
				fmt.Fprintf(&b, "limit_req_zone $binary_remote_addr zone=%s:10m rate=%s;\n", name, rate)
				seen[name] = true
			}
		}
	}
	if len(seen) == 0 {
		b.WriteString("# no active zones\n")
	}
	return b.String()
}

func zoneName(site, kind string) string {
	safe := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			return r
		default:
			return '_'
		}
	}, sanitizeSite(site))
	return fmt.Sprintf("nu_%s_%s", kind, safe)
}

func nonzero(v, fallback int) int {
	if v <= 0 {
		return fallback
	}
	return v
}

func regexpQuote(s string) string {
	replacer := strings.NewReplacer(
		`.`, `\.`, `+`, `\+`, `?`, `\?`, `*`, `\*`,
		`(`, `\(`, `)`, `\)`, `[`, `\[`, `]`, `\]`,
		`{`, `\{`, `}`, `\}`, `|`, `\|`, `^`, `\^`, `$`, `\$`,
	)
	return replacer.Replace(s)
}
