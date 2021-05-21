package http

import (
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	*gin.Engine
}

func (server *HttpServer) setRouter() {
}

func NewHttpServer() *HttpServer {
	httpServer := &HttpServer{
		Engine: gin.Default(),
	}
	httpServer.setRouter()

	return httpServer
}
