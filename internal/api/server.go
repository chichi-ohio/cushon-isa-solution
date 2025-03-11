package api

import (
	"context"
	"fmt"
	"net/http"

	"cushion-isa/internal/config"
	"cushion-isa/internal/logger"

	"github.com/gin-gonic/gin"
)

type Server struct {
	server *http.Server
	router *gin.Engine
}

func NewServer(cfg config.ServerConfig, handler *Handler) *Server {
	router := gin.Default()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Serve static files
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	// Frontend route
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// API routes
	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/invest", handler.CreateInvestment)
		apiGroup.GET("/investments", handler.ListInvestments)
		apiGroup.GET("/investments/:id", handler.GetInvestment)
	}

	return &Server{
		router: router,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	logger.Infof("Server starting on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
