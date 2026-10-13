package nginx_log

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/0xJacky/Nginx-UI/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClearDomainLogRequiresPath(t *testing.T) {
	err := ClearDomainLog("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}

func TestClearDomainLogTruncatesWhitelistedFile(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "comics_access.log")
	require.NoError(t, os.WriteFile(logPath, []byte("line1\nline2\n"), 0o644))

	prev := settings.NginxSettings.LogDirWhiteList
	settings.NginxSettings.LogDirWhiteList = []string{dir}
	t.Cleanup(func() {
		settings.NginxSettings.LogDirWhiteList = prev
	})

	require.NoError(t, ClearDomainLog(logPath))

	data, err := os.ReadFile(logPath)
	require.NoError(t, err)
	assert.Empty(t, data)
}
