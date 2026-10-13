package media

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/0xJacky/Nginx-UI/internal/nginx"
)

const MaxUploadBytes = 50 << 20

func Mkdir(rootName, rel, name string) error {
	name = strings.TrimSpace(name)
	if name == "" || strings.Contains(name, "..") || strings.ContainsAny(name, `/\`) {
		return fmt.Errorf("invalid folder name")
	}
	parent, err := Resolve(rootName, rel)
	if err != nil {
		return err
	}
	target := filepath.Join(parent, name)
	if !isUnderMediaWhitelist(target) {
		return fmt.Errorf("path is not under media whitelist")
	}
	return nginx.MkdirAll(target, 0o755)
}

func DeletePaths(rootName string, relPaths []string) (int, error) {
	if len(relPaths) == 0 {
		return 0, fmt.Errorf("no paths provided")
	}
	deleted := 0
	for _, rel := range relPaths {
		abs, err := Resolve(rootName, rel)
		if err != nil {
			return deleted, err
		}
		base, err := findRootPath(rootName)
		if err != nil {
			return deleted, err
		}
		if abs == base {
			return deleted, fmt.Errorf("refusing to delete media root")
		}
		if err := nginx.RemoveAll(abs); err != nil {
			return deleted, err
		}
		deleted++
	}
	return deleted, nil
}

func SaveUpload(rootName, rel, filename string, reader io.Reader, size int64) error {
	filename = filepath.Base(strings.TrimSpace(filename))
	if filename == "" || filename == "." || filename == ".." {
		return fmt.Errorf("invalid filename")
	}
	if !IsImageExt(filename) {
		return fmt.Errorf("only image files are allowed")
	}
	if size > MaxUploadBytes {
		return fmt.Errorf("file exceeds %d MB limit", MaxUploadBytes>>20)
	}
	dir, err := Resolve(rootName, rel)
	if err != nil {
		return err
	}
	target := filepath.Join(dir, filename)
	if !isUnderMediaWhitelist(target) {
		return fmt.Errorf("path is not under media whitelist")
	}

	tmp := target + ".upload.tmp"
	f, err := nginx.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	limited := io.LimitReader(reader, MaxUploadBytes+1)
	written, copyErr := io.Copy(f, limited)
	closeErr := f.Close()
	if copyErr != nil {
		_ = nginx.Remove(tmp)
		return copyErr
	}
	if closeErr != nil {
		_ = nginx.Remove(tmp)
		return closeErr
	}
	if written > MaxUploadBytes {
		_ = nginx.Remove(tmp)
		return fmt.Errorf("file exceeds %d MB limit", MaxUploadBytes>>20)
	}
	if err := nginx.Rename(tmp, target); err != nil {
		_ = nginx.Remove(tmp)
		return err
	}
	return nil
}
