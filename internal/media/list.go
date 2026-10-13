package media

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/0xJacky/Nginx-UI/internal/nginx"
)

var imageExts = map[string]string{
	".webp": "image/webp",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".svg":  "image/svg+xml",
	".bmp":  "image/bmp",
	".avif": "image/avif",
}

func IsImageExt(name string) bool {
	_, ok := imageExts[strings.ToLower(filepath.Ext(name))]
	return ok
}

func MimeForName(name string) string {
	if mime, ok := imageExts[strings.ToLower(filepath.Ext(name))]; ok {
		return mime
	}
	return "application/octet-stream"
}

func List(rootName, rel string) ([]Entry, error) {
	abs, err := Resolve(rootName, rel)
	if err != nil {
		return nil, err
	}

	entries, err := nginx.ReadDir(abs)
	if err != nil {
		return nil, err
	}

	rel = strings.Trim(filepath.ToSlash(rel), "/")
	result := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if name == "" || strings.HasPrefix(name, ".") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		childRel := name
		if rel != "" {
			childRel = path.Join(rel, name)
		}
		ext := strings.ToLower(filepath.Ext(name))
		item := Entry{
			Name:      name,
			Path:      childRel,
			IsDir:     entry.IsDir(),
			Size:      info.Size(),
			ModTime:   info.ModTime().Unix(),
			Extension: ext,
		}
		if !item.IsDir {
			item.MimeType = MimeForName(name)
			item.IsImage = IsImageExt(name)
		}
		result = append(result, item)
	}
	return result, nil
}
