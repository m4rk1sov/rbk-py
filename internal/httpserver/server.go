package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m4rk1sov/rbk-py/internal/config"
)

const apiVersion = "/api/v1"

type Server struct {
	cfg    config.Config
	engine *gin.Engine
	http   *http.Server
}

func New(cfg *config.Config) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	api := r.Group(filepath.Join(cfg.ServiceContextURL, apiVersion))
	api.Use(middleware.Auth(cfg.StaticToken, cfg.JWT.Secret))

	api.GET("/templates", func(c *gin.Context) {
		list, err := templates.List(cfg.TemplateDir)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err)
			return
		}
		web.OK(c, gin.H{"templates": list})
	})

	api.POST("/generate-html", handleGenerateHTML(cfg))
	api.POST("/generate-pdf", handleGeneratePDF(cfg))
	api.POST("/generate-docx", handleGenerateDOCX(cfg))
	api.POST("/generate-xlsx", handleGenerateXLSX(cfg))

	return &Server{
		cfg:    *cfg,
		engine: r,
		http: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
			Handler:      r,
			ReadTimeout:  cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
			IdleTimeout:  cfg.HTTP.IdleTimeout,
		},
	}
}

func (s *Server) Start() error {
	fmt.Printf("listening on: %s\n", s.http.Addr)
	return s.http.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.http.Shutdown(ctxTimeout)
}
