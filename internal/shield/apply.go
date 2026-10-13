package shield

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/0xJacky/Nginx-UI/internal/nginx"
	"github.com/0xJacky/Nginx-UI/internal/site"
	"github.com/0xJacky/Nginx-UI/model"
	"github.com/0xJacky/Nginx-UI/query"
)

const (
	markerStart = "# nginx-ui-shield-start"
	markerEnd   = "# nginx-ui-shield-end"
	zonesStart  = "# nginx-ui-shield-zones-start"
	zonesEnd    = "# nginx-ui-shield-zones-end"
)

func IsApplied(siteName string) bool {
	path, err := site.ResolveAvailablePath(siteName)
	if err != nil {
		return false
	}
	content, err := nginx.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), IncludePath(siteName))
}

func WriteRendered(cfg Config) error {
	if err := nginx.MkdirAll(Dir(), 0o755); err != nil {
		return err
	}
	body := Render(cfg)
	if err := nginx.WriteFile(SiteConfPath(cfg.Site), []byte(body), 0o644); err != nil {
		return err
	}
	return refreshZones()
}

func refreshZones() error {
	entries, err := nginx.ReadDir(Dir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var configs []Config
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, StateExt) {
			continue
		}
		siteName := strings.TrimSuffix(name, StateExt)
		cfg, err := Load(siteName)
		if err != nil {
			continue
		}
		configs = append(configs, cfg)
	}
	return nginx.WriteFile(ZonesPath(), []byte(RenderZones(configs)), 0o644)
}

func EnsureZonesInclude() error {
	path := nginx.GetConfEntryPath()
	contentBytes, err := nginx.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(contentBytes)
	includeLine := fmt.Sprintf("    include %s;", ZonesIncludePath())
	block := zonesStart + "\n" + includeLine + "\n" + zonesEnd

	if strings.Contains(content, ZonesIncludePath()) {
		return nil
	}

	updated, ok := injectIntoHTTP(content, block)
	if !ok {
		return fmt.Errorf("could not find http {} block in nginx.conf to install rate-limit zones")
	}
	return nginx.WriteFile(path, []byte(updated), 0o644)
}

func injectIntoHTTP(content, block string) (string, bool) {
	idx := strings.Index(content, "http")
	for idx >= 0 {
		rest := content[idx:]
		brace := strings.Index(rest, "{")
		if brace < 0 {
			return content, false
		}
		between := strings.TrimSpace(rest[len("http"):brace])
		if between == "" || strings.HasPrefix(between, "#") {
			insertAt := idx + brace + 1
			return content[:insertAt] + "\n" + block + "\n" + content[insertAt:], true
		}
		next := strings.Index(content[idx+1:], "http")
		if next < 0 {
			break
		}
		idx = idx + 1 + next
	}
	return content, false
}

func Apply(ctx context.Context, cfg Config) error {
	if !ValidSiteName(cfg.Site) {
		return ErrInvalidSite
	}
	if err := Save(cfg); err != nil {
		return err
	}
	if err := WriteRendered(cfg); err != nil {
		return err
	}
	if cfg.RateLimit.Enabled || cfg.UnderAttack.Enabled {
		if err := EnsureZonesInclude(); err != nil {
			return err
		}
	}

	path, err := site.ResolveAvailablePath(cfg.Site)
	if err != nil {
		return err
	}
	raw, err := nginx.ReadFile(path)
	if err != nil {
		return err
	}
	updated := injectSiteInclude(string(raw), cfg.Site)
	nsID, syncIDs := siteMeta(cfg.Site)
	if err := site.Save(ctx, cfg.Site, updated, true, nsID, syncIDs, model.PostSyncActionReloadNginx); err != nil {
		return err
	}
	cfg.Applied = true
	return Save(cfg)
}

func Detach(ctx context.Context, siteName string) error {
	if !ValidSiteName(siteName) {
		return ErrInvalidSite
	}
	path, err := site.ResolveAvailablePath(siteName)
	if err != nil {
		return err
	}
	raw, err := nginx.ReadFile(path)
	if err != nil {
		return err
	}
	updated := removeSiteInclude(string(raw))
	nsID, syncIDs := siteMeta(siteName)
	if err := site.Save(ctx, siteName, updated, true, nsID, syncIDs, model.PostSyncActionReloadNginx); err != nil {
		return err
	}
	cfg, _ := Load(siteName)
	cfg.Applied = false
	return Save(cfg)
}

func injectSiteInclude(content, siteName string) string {
	content = removeSiteInclude(content)
	include := fmt.Sprintf("%s\ninclude %s;\n%s", markerStart, IncludePath(siteName), markerEnd)

	servers := splitServerBlocks(content)
	if len(servers) == 0 {
		return content + "\n" + include + "\n"
	}

	var b strings.Builder
	last := 0
	for _, s := range servers {
		b.WriteString(content[last:s.openBrace+1])
		b.WriteString("\n    ")
		b.WriteString(strings.ReplaceAll(include, "\n", "\n    "))
		b.WriteString("\n")
		b.WriteString(content[s.openBrace+1 : s.end])
		last = s.end
	}
	b.WriteString(content[last:])
	return b.String()
}

func removeSiteInclude(content string) string {
	for {
		start := strings.Index(content, markerStart)
		if start < 0 {
			break
		}
		end := strings.Index(content[start:], markerEnd)
		if end < 0 {
			break
		}
		end = start + end + len(markerEnd)
		for end < len(content) && (content[end] == '\n' || content[end] == '\r') {
			end++
		}
		content = content[:start] + content[end:]
	}
	return content
}

type serverSpan struct {
	openBrace int
	end       int
}

func splitServerBlocks(content string) []serverSpan {
	var spans []serverSpan
	search := 0
	for {
		idx := strings.Index(content[search:], "server")
		if idx < 0 {
			break
		}
		idx += search
		rest := content[idx:]
		braceRel := strings.Index(rest, "{")
		if braceRel < 0 {
			break
		}
		between := strings.TrimSpace(rest[len("server"):braceRel])
		if between != "" && !strings.HasPrefix(between, "#") {
			search = idx + len("server")
			continue
		}
		open := idx + braceRel
		end := findMatchingBrace(content, open)
		if end < 0 {
			break
		}
		spans = append(spans, serverSpan{openBrace: open, end: end + 1})
		search = end + 1
	}
	return spans
}

func findMatchingBrace(content string, open int) int {
	depth := 0
	for i := open; i < len(content); i++ {
		switch content[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func siteMeta(siteName string) (uint64, []uint64) {
	path, err := site.ResolveAvailablePath(siteName)
	if err != nil {
		return 0, nil
	}
	s := query.Site
	record, err := s.Where(s.Path.Eq(path)).First()
	if err != nil || record == nil {
		return 0, nil
	}
	return record.NamespaceID, record.SyncNodeIDs
}
