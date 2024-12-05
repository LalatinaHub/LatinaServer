package main

import (
	"context"
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
	"github.com/go-co-op/gocron"
)

var (
	loc, _ = time.LoadLocation("Asia/Jakarta")
)

func updateUsersQuota() {
	defer helper.CatchError(true)

	var isAnyExceed bool
	for _, user := range config.ReadSingConfig(CS.SingActiveConfigPath).Experimental.V2RayAPI.Stats.Users {
		if !db.UpdatePremiumQuota(user) {
			isAnyExceed = true
		}

		time.Sleep(500 * time.Millisecond)
	}

	if isAnyExceed {
		helper.ReloadService([]string{CS.ServiceLatinaServer}...)
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

	// Start async funtions
	s.StartAsync()

	for {
		relay.GatherRelays()
		runtimeDebug.FreeOSMemory()

		cancelCtx, cancel := context.WithCancel(context.Background())
		go RunCaddyWithContext(cancelCtx)
		go RunSingBoxWithContext(cancelCtx)
		go web.RunWebServiceWithContext(cancelCtx)

		for {
			osSignal := <-c
			fmt.Println("Stopping services...")
			time.Sleep(2 * time.Second)

			if osSignal == syscall.SIGHUP {
				cancel()
				break
			}

			return
		}
	}
}
