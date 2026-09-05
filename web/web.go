package web

import (
	"net/http"

	serviceWeb "github.com/LalatinaHub/LatinaServer/internal/service/web"
)

func WebServer() http.Handler {
	server := serviceWeb.NewServer()
	return server.SetupRouter()
}

