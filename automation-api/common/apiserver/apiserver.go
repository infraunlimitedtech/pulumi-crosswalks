package apiserver

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	listener net.Listener
	logger   *zap.Logger
	http     *http.Server
}

type Route struct {
	Path    string
	Handler gin.HandlerFunc
}

func New(addr string, logger *zap.Logger, handlers []Route) (*Server, error) {
	// Using net listen since it can assign any available port
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("create tcp listener: %w", err)
	}

	binded := l.Addr().(*net.TCPAddr).String()
	logger.Info(fmt.Sprintf("listen on %s", binded))

	// Disable color since it is pretty ugly :(
	gin.DisableConsoleColor()

	// Dirty hack for gin debug mode
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error":  "Not implemented",
			"status": "ERROR",
		})
	})

	router.GET("/readyz", readyz())

	for _, h := range handlers {
		router.GET(h.Path, h.Handler)
	}

	srv := &Server{
		listener: l,
		logger:   logger,
		http: &http.Server{
			Handler: router,
		},
	}

	return srv, nil
}

func (w *Server) Run() error {
	err := w.http.Serve(w.listener)
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve: %w", err)
	}

	return nil
}

func (w *Server) Addr() string {
	return w.listener.Addr().(*net.TCPAddr).String()
}

func (w *Server) Close() error {
	w.logger.Info("stopping web server")
	return w.http.Close()
}

func readyz() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	}
}
