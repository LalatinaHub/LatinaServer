package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	runtimeDebug "runtime/debug"
	"syscall"
	"time"

	"github.com/LalatinaHub/LatinaServer/config"
	"github.com/LalatinaHub/LatinaServer/config/relay"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/db"
	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/LalatinaHub/LatinaServer/web"
	"github.com/LalatinaHub/wstunnel/pkg/tunnel"
	"github.com/go-co-op/gocron"
	singbox "github.com/sagernet/sing-box/pkg/sing-box"
)

var (
	loc, _   = time.LoadLocation("Asia/Jakarta")
	wsTunnel = tunnel.Server{
		Host: "127.0.0.1",
		Port: CS.WSTunnelPort,
	}
)

func hotReload() {
	config.Write()
	helper.ReloadService([]string{CS.ServiceSingBox, CS.ServiceOpenresty}...)
}

func updateUsersQuota() {
	defer helper.CatchError(true)

	var isAnyExceed bool
	for _, user := range config.ReadSingConfig().Experimental.V2RayAPI.Stats.Users {
		if !db.UpdatePremiumQuota(user) {
			isAnyExceed = true
		}

		time.Sleep(500 * time.Millisecond)
	}

	if isAnyExceed {
		hotReload()
	}
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)

	s := gocron.NewScheduler(loc)
	s.Every(1).Day().At("00:00").Tag("hot-reload").Do(hotReload)
	s.Every(30).Minutes().Tag("get-relays").Do(relay.GatherRelays)
	s.Every(5).Minutes().Tag("update-quota").Do(updateUsersQuota)
	s.Every(1).Day().At("03:00").Tag("reboot").Do(func() {
		if err := exec.Command("reboot").Run(); err != nil {
			fmt.Println("Failed to reboot the server !", err)
		}
	})

	runtimeDebug.FreeOSMemory()
	go web.StartWebService()
	go wsTunnel.Run()
	go singbox.RunWithOptions(config.GenerateSingConfig())

	// TODO
	// Handle reload (SIGHUP)
	for {
		osSignal := <-c
		fmt.Println("Stopping services...")
		time.Sleep(2 * time.Second)

		if osSignal != nil {
			s.Stop()
			break
		}
	}
}
