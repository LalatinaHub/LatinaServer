package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	CS "github.com/LalatinaHub/LatinaServer/constant"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func RunWebServiceWithContext(ctx context.Context) error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", CS.WebServerPort),
		Handler:      h2c.NewHandler(WebServer(), &http2.Server{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	defer server.Close()

	go func() {
		log.Println("Starting gin...")
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}

		log.Println("Gin started!")
	}()

	<-ctx.Done()

	log.Printf("Gin stopped due to %s\n", ctx.Err())
	return ctx.Err()
}
