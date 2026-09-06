package web

import (
	"context"
	"net/http"
	"net/http/pprof"
	"os"
	"path/filepath"
	"time"

	"github.com/LalatinaHub/LatinaServer/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	healthHandler       *HealthHandler
	statusHandler       *StatusHandler
	proxyHandler        *ProxyHandler
	adminHandler        *AdminHandler
	trialHandler        *TrialHandler
	dashboardHandler    *DashboardHandler
	userSettingsHandler *UserSettingsHandler
	portalHandler       *PortalHandler
	indexHandler        *IndexHandler
	rateLimiter         *RateLimiter
}

func NewServer() *Server {
	return NewServerWithContext(context.Background())
}

func NewServerWithContext(ctx context.Context) *Server {
	return &Server{
		healthHandler:       NewHealthHandler(),
		statusHandler:       NewStatusHandler(),
		proxyHandler:        NewProxyHandler(),
		adminHandler:        NewAdminHandler(),
		trialHandler:        NewTrialHandler(),
		dashboardHandler:    NewDashboardHandler(),
		userSettingsHandler: NewUserSettingsHandler(),
		portalHandler:       NewPortalHandler(),
		indexHandler:        NewIndexHandler(),
		rateLimiter:         NewRateLimiterWithContext(ctx, 100, 1*time.Minute),
	}
}

func (s *Server) SetupRouter() http.Handler {
	r := gin.New()

	// Global Middlewares
	r.Use(gin.Recovery())
	r.Use(RequestLoggerMiddleware())
	r.Use(ErrorMiddleware())
	r.Use(GzipMiddleware())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Pprof endpoints for profiling
	pprofGroup := r.Group("/debug/pprof")
	{
		pprofGroup.GET("/", gin.WrapF(pprof.Index))
		pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
		pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
		pprofGroup.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		pprofGroup.GET("/block", gin.WrapH(pprof.Handler("block")))
		pprofGroup.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		pprofGroup.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		pprofGroup.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		pprofGroup.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}

	// Routes
	apiV1 := r.Group("/api/v1")
	apiV1.Use(s.rateLimiter.Middleware())
	{
		// Health endpoints
		apiV1.GET("/ping", s.healthHandler.Ping)
		apiV1.GET("/info", s.healthHandler.Info)

		// Status endpoint
		apiV1.GET("/status", s.statusHandler.Status)

		// Proxy endpoints
		apiV1.GET("/check", s.proxyHandler.Check)
		apiV1.GET("/relay", s.proxyHandler.Relays)

		// Trial endpoints (VMess 2 Mbps speed limit)
		apiV1.POST("/trial", s.trialHandler.GenerateTrial)
		apiV1.GET("/trial", s.trialHandler.GenerateTrial)

		// User settings endpoints (AdBlock toggle, preference management)
		apiV1.GET("/user/settings", s.userSettingsHandler.GetSettings)
		apiV1.POST("/user/settings", s.userSettingsHandler.UpdateSettings)

		// Self-service portal endpoints (Parity with LatinaBot)
		portalGroup := apiV1.Group("/portal")
		{
			portalGroup.GET("/profile", s.portalHandler.GetProfile)
			portalGroup.POST("/reset-uuid", s.portalHandler.ResetUUID)
			portalGroup.POST("/change-token", s.portalHandler.ChangeToken)
			portalGroup.POST("/update-config", s.portalHandler.UpdateConfig)
			portalGroup.POST("/toggle-adblock", s.portalHandler.ToggleAdblock)
			portalGroup.GET("/servers", s.portalHandler.GetServers)
			portalGroup.GET("/relays", s.portalHandler.GetRelays)
			portalGroup.GET("/status", s.portalHandler.GetStatus)
			portalGroup.GET("/wildcards", s.portalHandler.GetWildcards)
			portalGroup.GET("/info", s.portalHandler.GetInfo)
		}

		// Admin endpoint (protected by Authorization: Bearer <token>)
		adminGroup := apiV1.Group("/admin")
		adminGroup.Use(s.adminHandler.AuthMiddleware())
		{
			adminGroup.POST("/reload", s.adminHandler.Reload)
			adminGroup.GET("/reload", s.adminHandler.Reload)
		}

		// Backward compatibility: Admin reload via token in path (deprecated)
		if token := s.adminHandler.GetAPIToken(); token != "" {
			apiV1.GET("/"+token, func(c *gin.Context) {
				logger.Warn().Msg("DEPRECATED: /api/v1/:token used for reload. Please use POST /api/v1/admin/reload with Bearer header")
				s.adminHandler.Reload(c)
			})
		}
	}

	// Locate static distribution directory (Hugo output)
	distDir := getDistDir()
	if fi, err := os.Stat(distDir); err == nil && fi.IsDir() {
		r.Static("/css", filepath.Join(distDir, "css"))
		r.Static("/js", filepath.Join(distDir, "js"))
		r.Static("/fonts", filepath.Join(distDir, "fonts"))
		r.Static("/images", filepath.Join(distDir, "images"))
		r.Static("/posts", filepath.Join(distDir, "posts"))
		r.Static("/records", filepath.Join(distDir, "records"))
		r.Static("/about", filepath.Join(distDir, "about"))
		r.StaticFile("/favicon.ico", filepath.Join(distDir, "favicon.ico"))
		r.StaticFile("/sitemap.xml", filepath.Join(distDir, "sitemap.xml"))
		r.StaticFile("/index.xml", filepath.Join(distDir, "index.xml"))
	}

	// Customer self-service portal route
	r.GET("/portal", s.portalHandler.ServePortal)
	r.GET("/portal/*any", s.portalHandler.ServePortal)

	// Dashboard route
	r.GET("/dashboard", s.dashboardHandler.ServeDashboard)

	// Root route: serves camouflage publication index or redirects to /portal?token=xxx
	r.GET("/", s.indexHandler.ServeIndex)

	return r
}

func getDistDir() string {
	if envPath := os.Getenv("LATINA_WEB_DIST"); envPath != "" {
		if fi, err := os.Stat(envPath); err == nil && fi.IsDir() {
			return envPath
		}
	}

	candidates := []string{
		filepath.Join("web", "dist"),
		filepath.Join("..", "web", "dist"),
		filepath.Join("..", "..", "web", "dist"),
		filepath.Join("..", "..", "..", "web", "dist"),
		"/var/www/mipa",
	}

	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			return c
		}
	}

	if wd, err := os.Getwd(); err == nil {
		dir := wd
		for i := 0; i < 6; i++ {
			testPath := filepath.Join(dir, "web", "dist")
			if fi, err := os.Stat(testPath); err == nil && fi.IsDir() {
				return testPath
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	return filepath.Join("web", "dist")
}

