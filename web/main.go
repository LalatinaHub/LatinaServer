package web

import (
	"fmt"
	"net/http"
	"time"

	CS "github.com/LalatinaHub/LatinaServer/constant"
	"golang.org/x/sync/errgroup"
)

var (
	g   errgroup.Group
	web = &http.Server{
		Addr:         fmt.Sprintf(":%d", CS.WebServerPort),
		Handler:      WebServer(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
)

func StartWebService() {
	g.Go(func() error {
		return web.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}
}
