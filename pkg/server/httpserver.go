package server

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Option func(*Server) error

type Server struct {
	httpAddr string
	engine   *gin.Engine

	shutdownTimeout time.Duration
}

func NewHttpServer(ctx context.Context, host string, port uint, shutdownTimeout time.Duration, options ...Option) (context.Context, Server) {

	srv := Server{
		engine:          gin.Default(),
		httpAddr:        fmt.Sprintf("%s:%d", host, port),
		shutdownTimeout: shutdownTimeout,
	}

	for _, option := range options {
		err := option(&srv)
		if err != nil {
			panic(err)
		}
	}

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")

	srv.engine.Use(cors.New(config))
	srv.engine.Use(otelgin.Middleware("api"))

	srv.engine.GET("/api/health", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})
	srv.engine.GET("/health", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})

	srv.engine.GET("/", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})

	return serverContext(ctx), srv
}

func (s *Server) Register(handler func(router *gin.RouterGroup)) {
	handler(s.engine.Group("/api"))
}

func (s Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.engine.ServeHTTP(w, req)
}

func (s *Server) Run(ctx context.Context) error {
	log.Println("Server running on", s.httpAddr)

	srv := &http.Server{
		Addr:    s.httpAddr,
		Handler: s.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("httpserver shut down", err)
		}
	}()

	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctxShutDown)
}

func serverContext(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
	}()

	return ctx
}
