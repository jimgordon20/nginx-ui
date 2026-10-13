package shield

import (
	"encoding/json"
	"os"

	"github.com/0xJacky/Nginx-UI/internal/nginx"
)

func Load(site string) (Config, error) {
	cfg := DefaultConfig(site)
	if !ValidSiteName(site) {
		return cfg, ErrInvalidSite
	}
	path := SiteStatePath(site)
	data, err := nginx.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	cfg.Site = site
	cfg.Preview = Render(cfg)
	cfg.Applied = IsApplied(site)
	return cfg, nil
}

func Save(cfg Config) error {
	if !ValidSiteName(cfg.Site) {
		return ErrInvalidSite
	}
	if err := nginx.MkdirAll(Dir(), 0o755); err != nil {
		return err
	}
	cfg.Preview = ""
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return nginx.WriteFile(SiteStatePath(cfg.Site), data, 0o644)
}

func LoadAlerts() ([]AlertRule, error) {
	data, err := nginx.ReadFile(AlertsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return []AlertRule{}, nil
		}
		return nil, err
	}
	var rules []AlertRule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}
	if rules == nil {
		rules = []AlertRule{}
	}
	return rules, nil
}

func SaveAlerts(rules []AlertRule) error {
	if err := nginx.MkdirAll(Dir(), 0o755); err != nil {
		return err
	}
	if rules == nil {
		rules = []AlertRule{}
	}
	data, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return err
	}
	return nginx.WriteFile(AlertsPath(), data, 0o644)
}
