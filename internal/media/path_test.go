package media

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/0xJacky/Nginx-UI/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveRejectsTraversal(t *testing.T) {
	original := append([]string{}, settings.NginxSettings.MediaDirWhiteList...)
	t.Cleanup(func() {
		settings.NginxSettings.MediaDirWhiteList = original
	})

	base := t.TempDir()
	rootName := "site"
	require.NoError(t, os.MkdirAll(filepath.Join(base, rootName, "images"), 0o755))
	settings.NginxSettings.MediaDirWhiteList = []string{base}

	_, err := Resolve(rootName, "../outside")
	assert.Error(t, err)

	_, err = Resolve("../"+rootName, "")
	assert.Error(t, err)

	_, err = Resolve(rootName+"/nested", "")
	assert.Error(t, err)

	abs, err := Resolve(rootName, "images/card.webp")
	require.NoError(t, err)
	assert.Equal(t, filepath.Clean(filepath.Join(base, rootName, "images", "card.webp")), abs)
}

func TestWhitelistDefaultsToVarWww(t *testing.T) {
	original := append([]string{}, settings.NginxSettings.MediaDirWhiteList...)
	t.Cleanup(func() {
		settings.NginxSettings.MediaDirWhiteList = original
	})

	settings.NginxSettings.MediaDirWhiteList = nil
	assert.Equal(t, []string{"/var/www"}, whitelist())
}

func TestIsImageExt(t *testing.T) {
	assert.True(t, IsImageExt("card.WEBP"))
	assert.True(t, IsImageExt("cover.png"))
	assert.False(t, IsImageExt("notes.txt"))
}
