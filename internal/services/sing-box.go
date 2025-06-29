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
	"github.com/sagernet/sing-box/include"
)

func RunSingBoxWithContext(ctx context.Context) error {
	if _, err := os.Stat(CS.SING_LOG_PATH); err == nil {
		os.Remove(CS.SING_LOG_PATH)
	}

	config.GenerateSingConfig()
	singCtx := box.Context(context.Background(), include.InboundRegistry(), include.OutboundRegistry(), include.EndpointRegistry(), include.DNSTransportRegistry(), include.ServiceRegistry())
	instance, err := box.New(box.Options{
		Context: singCtx,
		Options: config.ReadSingConfig(CS.SING_ACTIVE_CONFIG_PATH),
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
