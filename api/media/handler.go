package media

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/0xJacky/Nginx-UI/internal/media"
	"github.com/0xJacky/Nginx-UI/internal/nginx"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
)

func ListRoots(c *gin.Context) {
	roots, err := media.ListRoots()
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": roots})
}

func ListEntries(c *gin.Context) {
	root := c.Query("root")
	path := c.Query("path")
	entries, err := media.List(root, path)
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"root":    root,
		"path":    path,
		"entries": entries,
	})
}

func GetFile(c *gin.Context) {
	root := c.Query("root")
	path := c.Query("path")
	abs, err := media.Resolve(root, path)
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	if !media.IsImageExt(filepath.Base(abs)) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "only image files can be previewed"})
		return
	}
	f, err := nginx.Open(abs)
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	defer f.Close()

	info, err := nginx.Stat(abs)
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	c.Header("Content-Type", media.MimeForName(filepath.Base(abs)))
	c.Header("Content-Length", strconv.FormatInt(info.Size(), 10))
	c.Header("Cache-Control", "private, max-age=60")
	_, _ = io.Copy(c.Writer, f)
}

type mkdirRequest struct {
	Root string `json:"root" binding:"required"`
	Path string `json:"path"`
	Name string `json:"name" binding:"required"`
}

func Mkdir(c *gin.Context) {
	var req mkdirRequest
	if !cosy.BindAndValid(c, &req) {
		return
	}
	if err := media.Mkdir(req.Root, req.Path, req.Name); err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "folder created"})
}

type deleteRequest struct {
	Root  string   `json:"root" binding:"required"`
	Paths []string `json:"paths" binding:"required"`
}

func Delete(c *gin.Context) {
	var req deleteRequest
	if !cosy.BindAndValid(c, &req) {
		return
	}
	count, err := media.DeletePaths(req.Root, req.Paths)
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "deleted",
		"count":   count,
	})
}

func Upload(c *gin.Context) {
	root := c.PostForm("root")
	path := c.PostForm("path")
	form, err := c.MultipartForm()
	if err != nil {
		cosy.ErrHandler(c, err)
		return
	}
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no files uploaded"})
		return
	}

	uploaded := 0
	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			cosy.ErrHandler(c, err)
			return
		}
		err = media.SaveUpload(root, path, header.Filename, file, header.Size)
		_ = file.Close()
		if err != nil {
			cosy.ErrHandler(c, err)
			return
		}
		uploaded++
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "uploaded",
		"uploaded": uploaded,
	})
}
