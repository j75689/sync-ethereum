package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	httpServer *http.Server
	engine     *gin.Engine
}

func (server *HttpServer) setRouter() {
}

func NewHttpServer() *HttpServer {
	httpServer := &HttpServer{
		engine: gin.Default(),
	}
	httpServer.setRouter()

	return httpServer
}

func (server *HttpServer) Run(addr string) error {
	server.httpServer = &http.Server{
		Addr:    addr,
		Handler: server.engine,
	}
	return server.httpServer.ListenAndServe()
}

func (server *HttpServer) Shutdown() error {
	return server.httpServer.Close()
}
