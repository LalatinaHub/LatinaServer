package services

import (
	"context"
	"log"
	"os"

	"github.com/LalatinaHub/LatinaServer/config"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/experimental"
	"github.com/sagernet/sing-box/experimental/clashapi"
	"github.com/sagernet/sing-box/experimental/v2rayapi"
)

func RunSingBoxWithContext(ctx context.Context) error {
	if _, err := os.Stat(CS.SingLogPath); err == nil {
		os.Remove(CS.SingLogPath)
	}

	config.GenerateSingConfig()
	instance, err := box.New(box.Options{
		Context: context.Background(),
		Options: config.ReadSingConfig(CS.SingActiveConfigPath),
	})
	if err != nil {
		panic(err)
	}

	defer instance.Close()

	experimental.RegisterClashServerConstructor(clashapi.NewServer)
	experimental.RegisterV2RayServerConstructor(v2rayapi.NewServer)
	go func() {
		log.Println("Starting sing-box...")

		if err = instance.Start(); err != nil {
			panic(err)
		}

		log.Println("sing-box started!")
	}()

	<-ctx.Done()

	log.Printf("sing-box stopped due to %s\n", ctx.Err())
	return ctx.Err()
}
