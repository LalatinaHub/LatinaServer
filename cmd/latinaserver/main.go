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
	"github.com/LalatinaHub/LatinaServer/internal/services"
	"github.com/go-co-op/gocron"
)

var (
	loc, _ = time.LoadLocation("Asia/Jakarta")
)

func updateUsersQuota() {
	defer helper.CatchError(true)

	var isRunningOutQuota bool
	for _, user := range config.ReadSingConfig(CS.SING_ACTIVE_CONFIG_PATH).Experimental.V2RayAPI.Stats.Users {
		if db.UpdateAndCheckPremiumQuota(user) {
			isRunningOutQuota = true
		}

		time.Sleep(100 * time.Millisecond)
	}

	if isRunningOutQuota {
		helper.ReloadService([]string{CS.SERVICE_LATINASERVER}...)
	}
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)

	s := gocron.NewScheduler(loc)
	s.Every(5).Minutes().Tag("update-quota").Do(updateUsersQuota)
	s.Every(1).Day().At("03:00").Tag("reboot").Do(func() {
		if err := exec.Command("reboot").Run(); err != nil {
			fmt.Println("Failed to reboot the server !", err)
		}
	})

	// Start async funtions
	s.StartAsync()

	relay.GatherRelays()
	for {
		runtimeDebug.FreeOSMemory()

		cancelCtx, cancel := context.WithCancel(context.Background())
		go services.RunCaddyWithContext(cancelCtx)
		go services.RunSingBoxWithContext(cancelCtx)
		go services.RunGinWithContext(cancelCtx)

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
