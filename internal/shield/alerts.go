package shield

import (
	"context"
	"fmt"
	"time"

	"github.com/0xJacky/Nginx-UI/internal/nginx_log"
	"github.com/0xJacky/Nginx-UI/internal/nginx_log/analytics"
	"github.com/0xJacky/Nginx-UI/internal/notification"
	"github.com/google/uuid"
)

func UpsertAlert(rule AlertRule) ([]AlertRule, error) {
	rules, err := LoadAlerts()
	if err != nil {
		return nil, err
	}
	if rule.ID == "" {
		rule.ID = uuid.NewString()
	}
	if rule.WindowMinutes <= 0 {
		rule.WindowMinutes = 60
	}
	found := false
	for i := range rules {
		if rules[i].ID == rule.ID {
			rules[i] = rule
			found = true
			break
		}
	}
	if !found {
		rules = append(rules, rule)
	}
	if err := SaveAlerts(rules); err != nil {
		return nil, err
	}
	return rules, nil
}

func DeleteAlert(id string) ([]AlertRule, error) {
	rules, err := LoadAlerts()
	if err != nil {
		return nil, err
	}
	out := make([]AlertRule, 0, len(rules))
	for _, rule := range rules {
		if rule.ID != id {
			out = append(out, rule)
		}
	}
	if err := SaveAlerts(out); err != nil {
		return nil, err
	}
	return out, nil
}

func EvaluateAlerts(ctx context.Context) ([]AlertEvaluation, error) {
	rules, err := LoadAlerts()
	if err != nil {
		return nil, err
	}
	svc := nginx_log.GetAnalytics()
	if svc == nil {
		return nil, fmt.Errorf("analytics service unavailable")
	}

	now := time.Now().Unix()
	results := make([]AlertEvaluation, 0, len(rules))
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		window := int64(rule.WindowMinutes) * 60
		if window <= 0 {
			window = 3600
		}
		req := &analytics.DomainAnalyticsRequest{
			LogPath:   rule.LogPath,
			StartTime: now - window,
			EndTime:   now,
		}
		data, err := svc.GetDomainAnalytics(ctx, req)
		if err != nil {
			results = append(results, AlertEvaluation{
				Rule:      rule,
				Triggered: false,
				Message:   err.Error(),
			})
			continue
		}

		value := metricValue(rule.Metric, data)
		triggered := value >= rule.Threshold
		msg := fmt.Sprintf("%s = %.0f (threshold %.0f)", rule.Metric, value, rule.Threshold)
		eval := AlertEvaluation{
			Rule:      rule,
			Value:     value,
			Triggered: triggered,
			Message:   msg,
		}
		if triggered && rule.NotifyOnTrigger {
			notification.Warning(rule.Name, msg, nil)
		}
		results = append(results, eval)
	}
	return results, nil
}

func metricValue(metric string, data *analytics.DomainAnalytics) float64 {
	if data == nil {
		return 0
	}
	switch metric {
	case "bandwidth":
		return float64(data.Summary.TotalTraffic)
	case "requests":
		return float64(data.Summary.TotalPV)
	case "blocked":
		total := 0
		for _, item := range data.TopBlockedIPs {
			total += item.Count
		}
		return float64(total)
	case "server_errors":
		return data.ServerErrorRate
	case "client_errors":
		return data.ClientErrorRate
	default:
		return float64(data.Summary.TotalPV)
	}
}
