package web

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/LalatinaHub/LatinaServer/internal/config"
	"github.com/LalatinaHub/LatinaServer/internal/config/relay"
	"github.com/LalatinaHub/LatinaServer/internal/database"
	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/LalatinaHub/LatinaServer/internal/repository"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	"github.com/LalatinaHub/LatinaServer/pkg/systemctl"
	"github.com/LalatinaHub/LatinaServer/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PortalHandler handles the self-service web portal and its REST APIs.
type PortalHandler struct {
	repo            repository.UserSettingsRepository
	singConfigGen   func() error
	serviceReloader func() error
	customPath      string
	reloadMu        sync.Mutex
	reloadWg        sync.WaitGroup
}

// NewPortalHandler creates a new PortalHandler with database-backed repository.
func NewPortalHandler() *PortalHandler {
	db, err := database.GetDB()
	var repo repository.UserSettingsRepository
	if err == nil {
		repo = repository.NewUserSettingsRepository(db)
	}

	return &PortalHandler{
		repo:          repo,
		singConfigGen: config.GenerateSingConfig,
		serviceReloader: func() error {
			systemctl.Reload(config.ServiceLatinaServer)
			return nil
		},
	}
}

// NewPortalHandlerWithRepo creates a PortalHandler with custom dependencies for testing.
func NewPortalHandlerWithRepo(
	repo repository.UserSettingsRepository,
	singGen func() error,
	reload func() error,
	customPath string,
) *PortalHandler {
	return &PortalHandler{
		repo:            repo,
		singConfigGen:   singGen,
		serviceReloader: reload,
		customPath:      customPath,
	}
}

func (h *PortalHandler) getPortalPath() string {
	if h.customPath != "" {
		if _, err := os.Stat(h.customPath); err == nil {
			return h.customPath
		}
	}

	if envPath := os.Getenv("LATINA_PORTAL_PATH"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}
	}

	// Candidates for portal HTML
	candidates := []string{
		filepath.Join("web", "dist", "portal", "index.html"),
		filepath.Join("web", "dist", "portal.html"),
		filepath.Join("web", "portal.html"),
		filepath.Join("..", "web", "dist", "portal", "index.html"),
		filepath.Join("..", "web", "dist", "portal.html"),
		filepath.Join("..", "..", "web", "dist", "portal", "index.html"),
		filepath.Join("..", "..", "web", "dist", "portal.html"),
		filepath.Join("..", "..", "..", "web", "dist", "portal", "index.html"),
		filepath.Join("..", "..", "..", "web", "dist", "portal.html"),
		"/var/www/mipa/portal/index.html",
		"/var/www/mipa/portal.html",
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	// Search upward
	if wd, err := os.Getwd(); err == nil {
		dir := wd
		for i := 0; i < 6; i++ {
			pIndex := filepath.Join(dir, "web", "dist", "portal", "index.html")
			if _, err := os.Stat(pIndex); err == nil {
				return pIndex
			}
			p := filepath.Join(dir, "web", "dist", "portal.html")
			if _, err := os.Stat(p); err == nil {
				return p
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

// ServePortal serves the self-service portal single-page application.
func (h *PortalHandler) ServePortal(c *gin.Context) {
	path := h.getPortalPath()
	if path == "" {
		c.String(http.StatusNotFound, "Portal web file not found")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.File(path)
}

// ProfileResponse represents the payload sent to the portal.
type ProfileResponse struct {
	ID              int64          `json:"id"`
	Token           string         `json:"token"`
	Password        string         `json:"password"`
	Expired         string         `json:"expired"`
	ExpiredHuman    string         `json:"expired_human"`
	IsExpired       bool           `json:"is_expired"`
	IsDonator       bool           `json:"is_donator"`
	Quota           int64          `json:"quota"`
	QuotaHuman      string         `json:"quota_human"`
	ServerCode      string         `json:"server_code"`
	ServerDomain    string         `json:"server_domain"`
	Relay           string         `json:"relay"`
	VPN             string         `json:"vpn"`
	Adblock         bool           `json:"adblock"`
	SubURL          string         `json:"sub_url"`
	ClashURL        string         `json:"clash_url"`
	SingboxURL      string         `json:"singbox_url"`
	ShadowrocketURL string         `json:"shadowrocket_url"`
	RawURI          string         `json:"raw_uri"`
	Servers         []model.Server `json:"servers"`
	Relays          []string       `json:"relays"`
	Wildcards       []string       `json:"wildcards"`
}

// GetProfile returns the comprehensive user profile and subscription metadata.
func (h *PortalHandler) GetProfile(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required. Provide ?token=... or Authorization: Bearer <token>"})
		return
	}

	if h.repo == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database service unavailable"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	user, err := h.repo.GetUserByToken(ctx, token)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user: " + err.Error()})
		return
	}

	servers, _ := h.repo.GetServers(ctx)
	wildcards, _ := h.repo.GetWildcards(ctx)

	// Available relays
	relaysList := relay.GetRelays()
	relayCCMap := make(map[string]bool)
	for _, r := range relaysList {
		if r.CountryCode != "" {
			relayCCMap[strings.ToUpper(r.CountryCode)] = true
		}
	}
	var relayCCs []string
	for cc := range relayCCMap {
		relayCCs = append(relayCCs, cc)
	}
	sort.Strings(relayCCs)
	allRelays := append([]string{"Tanpa Relay"}, relayCCs...)

	// Find assigned server domain
	serverDomain := c.Request.Host
	for _, s := range servers {
		if s.Code == user.ServerCode && s.Domain != "" {
			serverDomain = s.Domain
			break
		}
	}

	now := time.Now()
	isExpired := now.After(user.Expired)
	isDonator := !isExpired

	var expiredHuman string
	if isExpired {
		expiredHuman = "Kedaluwarsa"
	} else {
		days := int(user.Expired.Sub(now).Hours() / 24)
		if days > 0 {
			expiredHuman = fmt.Sprintf("%d hari lagi", days)
		} else {
			hours := int(user.Expired.Sub(now).Hours())
			expiredHuman = fmt.Sprintf("%d jam lagi", hours)
		}
	}

	// URLs & Schemes
	proto := "https"
	if c.Request.TLS == nil && !strings.Contains(c.Request.Host, "localhost") {
		proto = "https"
	} else if strings.Contains(c.Request.Host, "localhost") {
		proto = "http"
	}
	baseURL := fmt.Sprintf("%s://%s", proto, c.Request.Host)
	subURL := fmt.Sprintf("%s/sub?token=%s", baseURL, user.Token)

	clashURL := fmt.Sprintf("clash://install-config?url=%s", url.QueryEscape(subURL))
	singboxURL := fmt.Sprintf("sing-box://import-remote-profile?url=%s", url.QueryEscape(subURL))
	shadowrocketURL := fmt.Sprintf("shadowrocket://add/sub://%s", base64.URLEncoding.EncodeToString([]byte(subURL)))

	// Generate raw URI for direct import
	rawURI := generateRawURI(user, serverDomain)

	resp := ProfileResponse{
		ID:              user.ID,
		Token:           user.Token,
		Password:        user.Password,
		Expired:         user.Expired.Format("2006-01-02"),
		ExpiredHuman:    expiredHuman,
		IsExpired:       isExpired,
		IsDonator:       isDonator,
		Quota:           user.Quota,
		QuotaHuman:      formatQuotaBytes(user.Quota),
		ServerCode:      user.ServerCode,
		ServerDomain:    serverDomain,
		Relay:           user.Relay,
		VPN:             user.VPN,
		Adblock:         user.Adblock,
		SubURL:          subURL,
		ClashURL:        clashURL,
		SingboxURL:      singboxURL,
		ShadowrocketURL: shadowrocketURL,
		RawURI:          rawURI,
		Servers:         servers,
		Relays:          allRelays,
		Wildcards:       wildcards,
	}

	c.JSON(http.StatusOK, resp)
}

// ResetUUID generates a fresh UUID for the user and reloads sing-box.
func (h *PortalHandler) ResetUUID(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		var req struct {
			Token string `json:"token"`
		}
		if err := c.ShouldBindJSON(&req); err == nil && req.Token != "" {
			token = req.Token
		}
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	if h.repo == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database service unavailable"})
		return
	}

	newUUID := uuid.New().String()
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.repo.UpdateUserUUID(ctx, token, newUUID); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update UUID: " + err.Error()})
		return
	}

	h.applyConfigAndReloadAsync()

	c.JSON(http.StatusOK, gin.H{
		"message":  "UUID successfully updated",
		"password": newUUID,
	})
}

// ChangeToken generates a fresh 8-char token for the user.
func (h *PortalHandler) ChangeToken(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		var req struct {
			Token string `json:"token"`
		}
		if err := c.ShouldBindJSON(&req); err == nil && req.Token != "" {
			token = req.Token
		}
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	if h.repo == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database service unavailable"})
		return
	}

	newToken := generateRandomToken(8)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.repo.UpdateUserToken(ctx, token, newToken); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update token: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token successfully updated",
		"token":   newToken,
	})
}

// UpdateConfigRequest represents the payload for changing VPN configuration.
type UpdateConfigRequest struct {
	Token      string `json:"token"`
	VPN        string `json:"vpn"`
	ServerCode string `json:"server_code"`
	Relay      string `json:"relay"`
}

// UpdateConfig updates VPN protocol, server code, and relay country.
func (h *PortalHandler) UpdateConfig(c *gin.Context) {
	var req UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	token := req.Token
	if token == "" {
		token = extractToken(c)
	}
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	req.VPN = strings.ToLower(strings.TrimSpace(req.VPN))
	if req.VPN != "vmess" && req.VPN != "vless" && req.VPN != "trojan" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid VPN protocol. Must be vmess, vless, or trojan"})
		return
	}

	if req.Relay == "Tanpa Relay" {
		req.Relay = ""
	}

	if h.repo == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database service unavailable"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.repo.UpdateUserConfig(ctx, token, req.VPN, req.ServerCode, req.Relay); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update config: " + err.Error()})
		return
	}

	h.applyConfigAndReloadAsync()

	c.JSON(http.StatusOK, gin.H{
		"message":     "VPN configuration updated successfully",
		"vpn":         req.VPN,
		"server_code": req.ServerCode,
		"relay":       req.Relay,
	})
}

// ToggleAdblock toggles the adblock state for the user.
func (h *PortalHandler) ToggleAdblock(c *gin.Context) {
	token := extractToken(c)
	var req struct {
		Token   string `json:"token"`
		Adblock *bool  `json:"adblock"`
	}
	if err := c.ShouldBindJSON(&req); err == nil {
		if req.Token != "" {
			token = req.Token
		}
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	if h.repo == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database service unavailable"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	user, err := h.repo.GetUserByToken(ctx, token)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user: " + err.Error()})
		return
	}

	targetState := !user.Adblock
	if req.Adblock != nil {
		targetState = *req.Adblock
	}

	if err := h.repo.UpdateUserAdblock(ctx, token, targetState); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle adblock: " + err.Error()})
		return
	}

	h.applyConfigAndReloadAsync()

	c.JSON(http.StatusOK, gin.H{
		"message": "AdBlock updated successfully",
		"adblock": targetState,
	})
}

// GetServers returns all edge servers with current load.
func (h *PortalHandler) GetServers(c *gin.Context) {
	if h.repo == nil {
		c.JSON(http.StatusOK, []model.Server{})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	servers, err := h.repo.GetServers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch servers: " + err.Error()})
		return
	}
	if servers == nil {
		servers = []model.Server{}
	}

	c.JSON(http.StatusOK, servers)
}

// GetRelays returns all available relay country codes.
func (h *PortalHandler) GetRelays(c *gin.Context) {
	relaysList := relay.GetRelays()
	relayCCMap := make(map[string]bool)
	for _, r := range relaysList {
		if r.CountryCode != "" {
			relayCCMap[strings.ToUpper(r.CountryCode)] = true
		}
	}
	var relayCCs []string
	for cc := range relayCCMap {
		relayCCs = append(relayCCs, cc)
	}
	sort.Strings(relayCCs)

	res := append([]string{"Tanpa Relay"}, relayCCs...)
	c.JSON(http.StatusOK, res)
}

// GetStatus returns the current hardware metrics and system telemetry.
func (h *PortalHandler) GetStatus(c *gin.Context) {
	status := util.GetServerStatus()
	c.JSON(http.StatusOK, status)
}

// GetWildcards returns all registered wildcard domains.
func (h *PortalHandler) GetWildcards(c *gin.Context) {
	if h.repo == nil {
		c.JSON(http.StatusOK, []string{})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	wildcards, _ := h.repo.GetWildcards(ctx)
	if wildcards == nil {
		wildcards = []string{}
	}
	c.JSON(http.StatusOK, wildcards)
}

// GetInfo returns community, donation, and disclaimer info.
func (h *PortalHandler) GetInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"telegram_group": "https://t.me/foolvpn",
		"website":        "https://foolvpn.web.id",
		"trakteer_url":   "https://trakteer.id/dickymuliafiqri/tip",
		"donation_info":  "1x donasi untuk 29 hari premium, berapapun jumlahnya.",
		"disclaimer":     "Semua akun yang disediakan oleh API merupakan akun yang tersedia secara bebas di Internet. Layanan ini tidak melakukan aktivitas tracing, logging, atau semacamnya!",
	})
}

func (h *PortalHandler) applyConfigAndReloadAsync() {
	h.reloadWg.Add(1)
	go func() {
		defer h.reloadWg.Done()
		h.reloadMu.Lock()
		defer h.reloadMu.Unlock()

		if h.singConfigGen != nil {
			if err := h.singConfigGen(); err != nil {
				logger.Error().Err(err).Msg("PortalHandler: failed to regenerate sing-box config")
				return
			}
		}

		if h.serviceReloader != nil {
			if err := h.serviceReloader(); err != nil {
				logger.Error().Err(err).Msg("PortalHandler: failed to reload service")
			}
		}
	}()
}

// WaitForAsyncReload blocks until all pending async reload operations complete.
// This is useful for testing and deterministic synchronization.
func (h *PortalHandler) WaitForAsyncReload() {
	h.reloadWg.Wait()
}

func formatQuotaBytes(b int64) string {
	if b <= 0 {
		return "0 MB"
	}
	const unit = 1024
	if b < unit*unit {
		return fmt.Sprintf("%.2f KB", float64(b)/float64(unit))
	}
	if b < unit*unit*unit {
		return fmt.Sprintf("%.2f MB", float64(b)/float64(unit*unit))
	}
	return fmt.Sprintf("%.2f GB", float64(b)/float64(unit*unit*unit))
}

func generateRandomToken(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	randBytes := make([]byte, n)
	_, _ = rand.Read(randBytes)
	for i := range b {
		b[i] = charset[int(randBytes[i])%len(charset)]
	}
	return string(b)
}

func generateRawURI(user *model.User, domain string) string {
	if user.Password == "" || domain == "" {
		return ""
	}

	switch strings.ToLower(user.VPN) {
	case "trojan":
		return fmt.Sprintf("trojan://%s@%s:443?security=tls&type=ws&path=%%2Ftrojan&sni=%s#%s",
			url.PathEscape(user.Password), domain, domain, url.QueryEscape(user.ServerCode))
	case "vless":
		return fmt.Sprintf("vless://%s@%s:443?path=%%2Fvless&security=tls&encryption=none&type=ws&sni=%s#%s",
			url.PathEscape(user.Password), domain, domain, url.QueryEscape(user.ServerCode))
	case "vmess":
		vmessConfig := map[string]any{
			"v":    "2",
			"ps":   user.ServerCode,
			"add":  domain,
			"port": "443",
			"id":   user.Password,
			"aid":  "0",
			"scy":  "auto",
			"net":  "ws",
			"type": "none",
			"host": domain,
			"path": "/vmess",
			"tls":  "tls",
			"sni":  domain,
		}
		raw, _ := json.Marshal(vmessConfig)
		return "vmess://" + base64.StdEncoding.EncodeToString(raw)
	default:
		return ""
	}
}
