package shield

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderHotlinkAndCacheCombined(t *testing.T) {
	cfg := ImageSitePreset("comics.conf")
	out := Render(cfg)
	assert.Contains(t, out, "valid_referers")
	assert.Contains(t, out, "limit_req zone=")
	assert.Contains(t, out, "Cache-Control")
	assert.Equal(t, 1, strings.Count(out, "location ~*"))
}

func TestRenderZones(t *testing.T) {
	cfg := ImageSitePreset("site1")
	cfg.UnderAttack.Enabled = true
	out := RenderZones([]Config{cfg})
	assert.Contains(t, out, "limit_req_zone")
	assert.Contains(t, out, "nu_media_")
	assert.Contains(t, out, "nu_attack_")
}
