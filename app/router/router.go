package router

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilhamfzri/pendek.in/config"
)

type Server struct {
	Server *http.Server
	Router *gin.Engine
}

func NewServer(cfg config.ServerConfig) *Server {
	router := gin.Default()

	readTimeout := time.Duration(cfg.ReadTimeout)
	writeTimeout := time.Duration(cfg.WriteTimeout)
	serverPort := fmt.Sprintf(":%s", strconv.Itoa(cfg.Port))

	server := &http.Server{
		Addr:         serverPort,
		Handler:      router,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
	}

	return &Server{
		Server: server,
		Router: router,
	}
}

func (server *Server) Run() {
	server.Server.ListenAndServe()
}
