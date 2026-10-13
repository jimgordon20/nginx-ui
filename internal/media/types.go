package media

type RootInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Size  int64  `json:"size"`
	Count int    `json:"count"`
}

type Entry struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	IsDir     bool   `json:"is_dir"`
	Size      int64  `json:"size"`
	ModTime   int64  `json:"mod_time"`
	MimeType  string `json:"mime_type,omitempty"`
	IsImage   bool   `json:"is_image"`
	Extension string `json:"extension,omitempty"`
}
