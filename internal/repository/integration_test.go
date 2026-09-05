//go:build integration

package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/LalatinaHub/LatinaServer/internal/database"
	"github.com/LalatinaHub/LatinaServer/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTursoIntegration(t *testing.T) {
	dbURL := os.Getenv("TURSO_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TURSO_DATABASE_URL environment variable not set; skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.GetDB()
	require.NoError(t, err)

	err = database.Ping(ctx)
	require.NoError(t, err)

	serverRepo := repository.NewServerRepository(db)
	servers, err := serverRepo.GetAll(ctx)
	assert.NoError(t, err)
	t.Logf("Fetched %d servers from live Turso database", len(servers))
}
