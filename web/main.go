package web

import (
	"fmt"
	"net/http"
	"time"

	CS "github.com/LalatinaHub/LatinaServer/constant"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func StartWebService() {
	g.Go(func() error {
		return (&http.Server{
			Addr:         fmt.Sprintf(":%d", CS.WebServerPort),
			Handler:      h2c.NewHandler(WebServer(), &http2.Server{}),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}).ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}
}
