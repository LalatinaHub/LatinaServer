package web

import (
	"context"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/LalatinaHub/LatinaServer/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	healthHandler *HealthHandler
	statusHandler *StatusHandler
	proxyHandler  *ProxyHandler
	adminHandler  *AdminHandler
	rateLimiter   *RateLimiter
}

func NewServer() *Server {
	return NewServerWithContext(context.Background())
}

func NewServerWithContext(ctx context.Context) *Server {
	return &Server{
		healthHandler: NewHealthHandler(),
		statusHandler: NewStatusHandler(),
		proxyHandler:  NewProxyHandler(),
		adminHandler:  NewAdminHandler(),
		rateLimiter:   NewRateLimiterWithContext(ctx, 100, 1*time.Minute),
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

	// Welcome route
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Gin!")
	})

	return r
}

