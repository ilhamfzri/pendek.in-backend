package router

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

	// .- Initialize Resources Directories
	os.MkdirAll(cfg.ResourcesDirPath, os.ModePerm)
	userProfilePicturePath := filepath.Join(cfg.ResourcesDirPath, "profile_pic")
	customThumbnailPicturePath := filepath.Join(cfg.ResourcesDirPath, "thumbnail")

	os.MkdirAll(userProfilePicturePath, os.ModePerm)
	os.MkdirAll(customThumbnailPicturePath, os.ModePerm)

	os.Setenv("PROFILE_IMG_DIR", userProfilePicturePath)
	os.Setenv("THUMBNAIL_IMG_DIR", customThumbnailPicturePath)

	return &Server{
		Server: server,
		Router: router,
	}
}

func (server *Server) Run() {
	server.Server.ListenAndServe()
}
