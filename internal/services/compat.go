package services

import (
	"context"

	"github.com/LalatinaHub/LatinaServer/internal/service/caddy"
	"github.com/LalatinaHub/LatinaServer/internal/service/singbox"
	"github.com/LalatinaHub/LatinaServer/internal/service/web"
)

func RunCaddyWithContext(ctx context.Context) error {
	return caddy.RunWithContext(ctx)
}

func RunSingBoxWithContext(ctx context.Context) error {
	return singbox.RunWithContext(ctx)
}

func RunGinWithContext(ctx context.Context) error {
	return web.RunWithContext(ctx)
}
