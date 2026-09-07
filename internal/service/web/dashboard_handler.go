package web

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	customPath string
}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func NewDashboardHandlerWithPath(customPath string) *DashboardHandler {
	return &DashboardHandler{customPath: customPath}
}

func (h *DashboardHandler) getDashboardPath() string {
	if h.customPath != "" {
		if _, err := os.Stat(h.customPath); err == nil {
			return h.customPath
		}
	}

	if envPath := os.Getenv("LATINA_DASHBOARD_PATH"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}
	}

	// 1. Production Linux path installed by script/install.sh
	if _, err := os.Stat("/var/www/mipa/dashboard.html"); err == nil {
		return "/var/www/mipa/dashboard.html"
	}

	// 2. Relative candidates
	candidates := []string{
		filepath.Join("web", "dist", "dashboard.html"),
		filepath.Join("web", "static", "dashboard.html"),
		filepath.Join("..", "web", "dist", "dashboard.html"),
		filepath.Join("..", "web", "static", "dashboard.html"),
		filepath.Join("..", "..", "web", "dist", "dashboard.html"),
		filepath.Join("..", "..", "web", "static", "dashboard.html"),
		filepath.Join("..", "..", "..", "web", "dist", "dashboard.html"),
		filepath.Join("..", "..", "..", "web", "static", "dashboard.html"),
		filepath.Join("..", "..", "..", "..", "web", "dist", "dashboard.html"),
		filepath.Join("..", "..", "..", "..", "web", "static", "dashboard.html"),
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
			testPath := filepath.Join(dir, "web", "dist", "dashboard.html")
			if _, err := os.Stat(testPath); err == nil {
				return testPath
			}
			staticPath := filepath.Join(dir, "web", "static", "dashboard.html")
			if _, err := os.Stat(staticPath); err == nil {
				return staticPath
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

func (h *DashboardHandler) ServeDashboard(c *gin.Context) {
	path := h.getDashboardPath()
	if path == "" {
		c.String(http.StatusNotFound, "Dashboard file not found")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.File(path)
}
