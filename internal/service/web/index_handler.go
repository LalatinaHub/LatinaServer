package web

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// IndexHandler serves the decoy camouflage research website.
type IndexHandler struct {
	customPath string
}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

func NewIndexHandlerWithPath(customPath string) *IndexHandler {
	return &IndexHandler{customPath: customPath}
}

func (h *IndexHandler) getIndexPath() string {
	if h.customPath != "" {
		if _, err := os.Stat(h.customPath); err == nil {
			return h.customPath
		}
		return ""
	}

	if envPath := os.Getenv("LATINA_INDEX_PATH"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}
	}

	// 1. Production Linux path installed by script/install.sh
	if _, err := os.Stat("/var/www/mipa/index.html"); err == nil {
		return "/var/www/mipa/index.html"
	}

	// 2. Relative candidates
	candidates := []string{
		filepath.Join("web", "dist", "index.html"),
		filepath.Join("..", "web", "dist", "index.html"),
		filepath.Join("..", "..", "web", "dist", "index.html"),
		filepath.Join("..", "..", "..", "web", "dist", "index.html"),
		filepath.Join("..", "..", "..", "..", "web", "dist", "index.html"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// 3. Search upward from working directory
	if wd, err := os.Getwd(); err == nil {
		dir := wd
		for i := 0; i < 6; i++ {
			testPath := filepath.Join(dir, "web", "dist", "index.html")
			if _, err := os.Stat(testPath); err == nil {
				return testPath
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	return ""
}

// ServeIndex serves the camouflage publication homepage or redirects if token is present.
func (h *IndexHandler) ServeIndex(c *gin.Context) {
	if token := c.Query("token"); token != "" {
		c.Redirect(http.StatusTemporaryRedirect, "/portal?token="+token)
		return
	}

	path := h.getIndexPath()
	if path == "" {
		c.String(http.StatusNotFound, "Index camouflage site not found")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.File(path)
}
