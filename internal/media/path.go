package media

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/0xJacky/Nginx-UI/internal/helper"
	"github.com/0xJacky/Nginx-UI/internal/nginx"
	"github.com/0xJacky/Nginx-UI/settings"
)

func whitelist() []string {
	list := append([]string{}, settings.NginxSettings.MediaDirWhiteList...)
	if len(list) == 0 {
		list = []string{"/var/www"}
	}
	return list
}

func isUnderMediaWhitelist(path string) bool {
	clean := filepath.Clean(path)
	for _, root := range whitelist() {
		if helper.IsUnderDirectory(clean, root) {
			return true
		}
	}
	return false
}

func Resolve(rootName, rel string) (abs string, err error) {
	rootName = strings.TrimSpace(rootName)
	if rootName == "" || strings.Contains(rootName, "..") || strings.ContainsAny(rootName, `/\`) {
		return "", fmt.Errorf("invalid media root")
	}

	base, err := findRootPath(rootName)
	if err != nil {
		return "", err
	}

	rel = strings.TrimSpace(rel)
	rel = strings.TrimPrefix(rel, "/")
	rel = filepath.ToSlash(rel)
	if strings.Contains(rel, "..") {
		return "", fmt.Errorf("invalid media path")
	}

	abs = filepath.Clean(filepath.Join(base, filepath.FromSlash(rel)))
	if !helper.IsUnderDirectory(abs, base) {
		return "", fmt.Errorf("path escapes media root")
	}
	if !isUnderMediaWhitelist(abs) {
		return "", fmt.Errorf("path is not under media whitelist")
	}
	return abs, nil
}

func findRootPath(rootName string) (string, error) {
	for _, base := range whitelist() {
		candidate := filepath.Clean(filepath.Join(base, rootName))
		if !helper.IsUnderDirectory(candidate, base) {
			continue
		}
		info, err := nginx.Stat(candidate)
		if err != nil {
			continue
		}
		if info.IsDir() {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("media root not found")
}

func ListRoots() ([]RootInfo, error) {
	seen := map[string]bool{}
	var roots []RootInfo
	for _, base := range whitelist() {
		entries, err := nginx.ReadDir(base)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			name := entry.Name()
			if name == "" || strings.HasPrefix(name, ".") {
				continue
			}
			if seen[name] {
				continue
			}
			seen[name] = true
			full := filepath.Join(base, name)
			var size int64
			var count int
			_ = walkCount(full, &size, &count)
			roots = append(roots, RootInfo{
				Name:  name,
				Path:  full,
				Size:  size,
				Count: count,
			})
		}
	}
	return roots, nil
}

func walkCount(dir string, size *int64, count *int) error {
	entries, err := nginx.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		full := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			_ = walkCount(full, size, count)
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		*count++
		*size += info.Size()
	}
	return nil
}
