package v2ray

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/LalatinaHub/LatinaServer/pkg/constant"
	appErrors "github.com/LalatinaHub/LatinaServer/pkg/errors"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	"github.com/v2fly/v2ray-core/v5/app/stats/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	cachedConn   *grpc.ClientConn
	cachedClient command.StatsServiceClient
	connMu       sync.Mutex
)

// getTransportCredentials resolves TLS or insecure credentials based on environment.
func getTransportCredentials() credentials.TransportCredentials {
	if certFile := os.Getenv("V2RAY_API_CERT"); certFile != "" {
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err == nil {
			return creds
		}
		logger.Warn().Err(err).Str("cert", certFile).Msg("Failed to load V2Ray TLS cert file, falling back")
	}

	if os.Getenv("V2RAY_API_TLS") == "true" {
		return credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: os.Getenv("V2RAY_API_TLS_INSECURE") == "true",
		})
	}

	return insecure.NewCredentials()
}

// getClient returns a cached, reusable gRPC client or connects if disconnected/nil.
func getClient() (command.StatsServiceClient, error) {
	connMu.Lock()
	defer connMu.Unlock()

	if cachedConn != nil {
		state := cachedConn.GetState()
		if state == connectivity.Ready || state == connectivity.Idle {
			return cachedClient, nil
		}
		if state == connectivity.Shutdown {
			_ = cachedConn.Close()
			cachedConn = nil
			cachedClient = nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, constant.V2RayAPIAddress,
		grpc.WithTransportCredentials(getTransportCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, appErrors.NewNetworkError("failed to connect to v2ray stats gRPC", err)
	}

	cachedConn = conn
	cachedClient = command.NewStatsServiceClient(conn)
	return cachedClient, nil
}

// Close closes the cached V2Ray gRPC client connection if active.
func Close() error {
	connMu.Lock()
	defer connMu.Unlock()
	if cachedConn != nil {
		err := cachedConn.Close()
		cachedConn = nil
		cachedClient = nil
		return err
	}
	return nil
}

// GetUserStats retrieves downlink traffic stats for a single user.
func GetUserStats(name string) int64 {
	client, err := getClient()
	if err != nil {
		logger.Warn().Err(err).Str("user", name).Msg("Failed to connect to v2ray gRPC stats")
		return 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := client.GetStats(ctx, &command.GetStatsRequest{
		Name:   fmt.Sprintf("user>>>%s>>>traffic>>>downlink", name),
		Reset_: true,
	})
	if err != nil {
		connMu.Lock()
		if cachedConn != nil {
			_ = cachedConn.Close()
			cachedConn = nil
			cachedClient = nil
		}
		connMu.Unlock()

		logger.Debug().Err(err).Str("user", name).Msg("Failed to get user stats from v2ray")
		return 0
	}

	if resp == nil || resp.Stat == nil {
		return 0
	}

	logger.Debug().
		Str("user", name).
		Int64("bytes", resp.Stat.Value).
		Int64("mb", resp.Stat.Value/1000000).
		Msg("V2Ray user traffic stats")

	return resp.Stat.Value
}

// GetUsersStatsBatch queries stats for multiple users reusing the single gRPC client connection.
func GetUsersStatsBatch(userNames []string) map[int64]int64 {
	usages := make(map[int64]int64)
	if len(userNames) == 0 {
		return usages
	}

	client, err := getClient()
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to connect to v2ray gRPC stats for batch")
		return usages
	}

	for _, userStr := range userNames {
		userID, err := strconv.ParseInt(userStr, 10, 64)
		if err != nil {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		resp, err := client.GetStats(ctx, &command.GetStatsRequest{
			Name:   fmt.Sprintf("user>>>%s>>>traffic>>>downlink", userStr),
			Reset_: true,
		})
		cancel()

		if err != nil {
			continue
		}

		if resp != nil && resp.Stat != nil && resp.Stat.Value > 0 {
			usages[userID] = resp.Stat.Value
		}
	}

	return usages
}

