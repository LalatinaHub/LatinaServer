package warp

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	"golang.org/x/crypto/curve25519"
)

const (
	DefaultWarpEndpoint   = "162.159.192.1:2408"
	DefaultPeerPublicKey  = "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo="
	DefaultIPv4           = "172.16.0.2/32"
	DefaultIPv6           = "2606:4700:110:8f81:85e8:bb2d:5c55:6834/128"
	DefaultConfigPath     = "/usr/local/etc/latinaserver/warp.json"
	CloudflareRegisterURL = "https://api.cloudflareclient.com/v0a3371/reg"
)

// Config represents a Cloudflare WARP WireGuard configuration.
type Config struct {
	PrivateKey    string   `json:"private_key"`
	PublicKey     string   `json:"public_key"`
	IPv4          string   `json:"ipv4"`
	IPv6          string   `json:"ipv6"`
	Reserved      [3]uint8 `json:"reserved"`
	Endpoint      string   `json:"endpoint"`
	PeerPublicKey string   `json:"peer_public_key"`
}

var (
	cachedConfig *Config
	configMu     sync.RWMutex
	customPath   string
)

// SetCustomConfigPath sets a custom storage path for warp.json (useful for testing).
func SetCustomConfigPath(path string) {
	configMu.Lock()
	defer configMu.Unlock()
	customPath = path
}

func getConfigPath() string {
	if customPath != "" {
		return customPath
	}
	return DefaultConfigPath
}

// GenerateKeyPair generates a Curve25519 keypair for WireGuard.
func GenerateKeyPair() (privateKeyBase64, publicKeyBase64 string, err error) {
	var privKey [32]byte
	if _, err := rand.Read(privKey[:]); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Clamp private key per WireGuard / RFC 7748 specification
	privKey[0] &= 248
	privKey[31] &= 127
	privKey[31] |= 64

	pubKey, err := curve25519.X25519(privKey[:], curve25519.Basepoint)
	if err != nil {
		return "", "", fmt.Errorf("curve25519 scalar multiplication failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(privKey[:]), base64.StdEncoding.EncodeToString(pubKey), nil
}

// LoadConfig loads cached WARP configuration from disk.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode warp config: %w", err)
	}

	if cfg.PrivateKey == "" || cfg.PeerPublicKey == "" {
		return nil, fmt.Errorf("invalid warp config content")
	}

	return &cfg, nil
}

// SaveConfig saves WARP configuration to disk.
func SaveConfig(path string, cfg *Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal warp config: %w", err)
	}

	return os.WriteFile(path, data, 0600)
}

// RegisterAccount registers a new WARP client account with Cloudflare API.
func RegisterAccount(ctx context.Context, client *http.Client) (*Config, error) {
	privKey, pubKey, err := GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	reqBody := map[string]any{
		"key":           pubKey,
		"install_id":    "",
		"fcm_token":     "",
		"tos":           time.Now().UTC().Format(time.RFC3339Nano),
		"model":         "PC",
		"serial_number": "",
		"locale":        "en_US",
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, CloudflareRegisterURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "okhttp/3.12.1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to cloudflare api failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cloudflare api returned status %d", resp.StatusCode)
	}

	var cfResp struct {
		Result struct {
			Config struct {
				ClientID string `json:"client_id"`
				Peers    []struct {
					PublicKey string `json:"public_key"`
					Endpoint  struct {
						V4   string `json:"v4"`
						Host string `json:"host"`
					} `json:"endpoint"`
				} `json:"peers"`
				Interface struct {
					Addresses struct {
						V4 string `json:"v4"`
						V6 string `json:"v6"`
					} `json:"addresses"`
				} `json:"interface"`
			} `json:"config"`
		} `json:"result"`
		Success bool `json:"success"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&cfResp); err != nil {
		return nil, fmt.Errorf("failed to parse cloudflare response: %w", err)
	}

	cfg := &Config{
		PrivateKey:    privKey,
		PublicKey:     pubKey,
		IPv4:          cfResp.Result.Config.Interface.Addresses.V4,
		IPv6:          cfResp.Result.Config.Interface.Addresses.V6,
		Endpoint:      DefaultWarpEndpoint,
		PeerPublicKey: DefaultPeerPublicKey,
		Reserved:      [3]uint8{0, 0, 0},
	}

	if cfg.IPv4 == "" {
		cfg.IPv4 = DefaultIPv4
	}
	if cfg.IPv6 == "" {
		cfg.IPv6 = DefaultIPv6
	}

	if len(cfResp.Result.Config.Peers) > 0 {
		peer := cfResp.Result.Config.Peers[0]
		if peer.PublicKey != "" {
			cfg.PeerPublicKey = peer.PublicKey
		}
		if peer.Endpoint.V4 != "" {
			cfg.Endpoint = peer.Endpoint.V4
		} else if peer.Endpoint.Host != "" {
			cfg.Endpoint = peer.Endpoint.Host
		}
	}

	if clientID := cfResp.Result.Config.ClientID; clientID != "" {
		if rawID, err := base64.StdEncoding.DecodeString(clientID); err == nil && len(rawID) >= 3 {
			cfg.Reserved = [3]uint8{rawID[0], rawID[1], rawID[2]}
		}
	}

	return cfg, nil
}

// GetOrInitWarpConfig retrieves or initializes WARP configuration.
func GetOrInitWarpConfig(ctx context.Context) (*Config, error) {
	configMu.Lock()
	defer configMu.Unlock()

	if cachedConfig != nil {
		return cachedConfig, nil
	}

	// 1. Environment Variable Overrides
	if envPriv := os.Getenv("WARP_PRIVATE_KEY"); envPriv != "" {
		endpoint := os.Getenv("WARP_ENDPOINT")
		if endpoint == "" {
			endpoint = DefaultWarpEndpoint
		}
		peerPub := os.Getenv("WARP_PEER_PUBLIC_KEY")
		if peerPub == "" {
			peerPub = DefaultPeerPublicKey
		}
		ipv4 := os.Getenv("WARP_IPV4")
		if ipv4 == "" {
			ipv4 = DefaultIPv4
		}
		ipv6 := os.Getenv("WARP_IPV6")
		if ipv6 == "" {
			ipv6 = DefaultIPv6
		}

		cachedConfig = &Config{
			PrivateKey:    envPriv,
			IPv4:          ipv4,
			IPv6:          ipv6,
			Endpoint:      endpoint,
			PeerPublicKey: peerPub,
			Reserved:      [3]uint8{0, 0, 0},
		}
		logger.Info().Msg("WARP configured from environment variables")
		return cachedConfig, nil
	}

	// 2. Load from disk cache
	cfgPath := getConfigPath()
	if cfg, err := LoadConfig(cfgPath); err == nil {
		cachedConfig = cfg
		logger.Info().Str("path", cfgPath).Msg("Loaded cached WARP configuration from disk")
		return cachedConfig, nil
	}

	// 3. Register a new account via Cloudflare Client API
	logger.Info().Msg("Registering new Cloudflare WARP account...")
	cfg, err := RegisterAccount(ctx, nil)
	if err != nil {
		logger.Warn().Err(err).Msg("Cloudflare API registration failed; using fallback keypair and defaults")
		privKey, pubKey, keyErr := GenerateKeyPair()
		if keyErr != nil {
			return nil, fmt.Errorf("failed to generate fallback keypair: %w", keyErr)
		}
		cfg = &Config{
			PrivateKey:    privKey,
			PublicKey:     pubKey,
			IPv4:          DefaultIPv4,
			IPv6:          DefaultIPv6,
			Endpoint:      DefaultWarpEndpoint,
			PeerPublicKey: DefaultPeerPublicKey,
			Reserved:      [3]uint8{0, 0, 0},
		}
	}

	// Persist to disk
	if err := SaveConfig(cfgPath, cfg); err != nil {
		logger.Warn().Err(err).Str("path", cfgPath).Msg("Failed to persist WARP config to disk")
	} else {
		logger.Info().Str("path", cfgPath).Msg("Saved new WARP configuration to disk")
	}

	cachedConfig = cfg
	return cachedConfig, nil
}
