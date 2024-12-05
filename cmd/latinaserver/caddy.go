package main

import (
	"context"
	"log"
	"os"

	"github.com/LalatinaHub/LatinaServer/config"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	caddy "github.com/caddyserver/caddy/v2"
)

func RunCaddyWithContext(ctx context.Context) error {
	if _, err := os.Stat(CS.CaddyLogPath); err == nil {
		os.Remove(CS.CaddyLogPath)
	}
	defer caddy.Stop()

	config.GenerateCaddyConfig()
	go func() {
		log.Println("Starting caddy...")
		if err := caddy.Run(config.ReadCaddyConfig(CS.CaddyActiveConfigPath)); err != nil {
			panic(err)
		}

		log.Println("Caddy started!")
	}()

	<-ctx.Done()

	log.Printf("Caddy stopped due to %s\n", ctx.Err())
	return ctx.Err()
}
