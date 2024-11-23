package main

import (
	"context"
	"fmt"
	"log"
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
	box "github.com/sagernet/sing-box"
)

var (
	loc, _   = time.LoadLocation("Asia/Jakarta")
	wsTunnel = tunnel.Server{
		Host: "127.0.0.1",
		Port: CS.WSTunnelPort,
	}
)

func hotReload() {
	helper.ReloadService([]string{CS.ServiceLatinaServer}...)
}

func startOpenresty() {
	config.WriteOpenrestyConfig()
	helper.ReloadService([]string{CS.ServiceOpenresty}...)
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
	s.Every(15).Minutes().Tag("update-quota").Do(updateUsersQuota)
	s.Every(1).Day().At("03:00").Tag("reboot").Do(func() {
		if err := exec.Command("reboot").Run(); err != nil {
			fmt.Println("Failed to reboot the server !", err)
		}
	})

	// Init
	relay.GatherRelays()

	// Start async funtions
	go web.StartWebService()
	s.StartAsync()

	for {
		var (
			options = config.GenerateSingConfig()
			quit    = make(chan bool)
		)

		os.Remove("/usr/local/etc/latinaserver/singbox.log")
		runtimeDebug.FreeOSMemory()
		startOpenresty()
		go func() {
			log.Println("Starting sing-box...")
			instance, err := box.New(box.Options{
				Context: context.Background(),
				Options: options,
			})
			if err != nil {
				panic(err)
			}

			if err = instance.Start(); err != nil {
				panic(err)
			}

			log.Println("sing-box started!")
			for {
				if <-quit {
					instance.Close()
					return
				}
			}
		}()

		for {
			osSignal := <-c
			fmt.Println("Stopping services...")
			time.Sleep(2 * time.Second)

			if osSignal == syscall.SIGHUP {
				quit <- true
				break
			}

			return
		}
	}
}
