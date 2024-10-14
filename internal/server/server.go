package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"session_manager/internal/api"
	"session_manager/internal/repository/pb/session_manager"
	"session_manager/internal/repository/postgres"
	"session_manager/internal/service"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
)

type Server interface {
	Run(ctx context.Context)
	Stop(ctx context.Context)
}

type server struct {
	router  *echo.Echo
	grpcCli *grpc.Server
}

func NewServer(env *Env) Server {
	s := server{
		router:  echo.New(),
		grpcCli: grpc.NewServer(),
	}

	// storage
	storage := postgres.NewStorage(env.pool)

	// service
	svc := service.New(storage)

	// service for grpc
	svcGrpc := service.NewServiceGrpc(storage)
	session_manager.RegisterSessionManagerServer(s.grpcCli, svcGrpc)

	// handlers
	hndl := api.NewHandlers(s.router.Logger, svc)

	// set middlewares
	s.router.Use(middleware.Logger(), middleware.Recover())

	// register handlers
	g := s.router.Group("/api/session-manager")
	g.GET("/dashboard", hndl.GetOnlineSessions)
	g.GET("/activity", hndl.GetUserActivity)

	sg := g.Group("/session")
	sg.POST("/on-campus", hndl.CreateSessionOnCampus)
	sg.POST("/on-platform", hndl.CreateSessionOnPlatform)

	return &s
}

func (s *server) Run(ctx context.Context) {
	ctxSignal, cancelSignal := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// start rest api server
	go func() {
		defer cancelSignal()

		if err := s.router.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Printf("server start error: %s\n", err.Error())
		}
	}()

	// start grpc server
	go func() {
		defer cancelSignal()
		addr := ":8181"

		log.Printf("grpc starts on port: %s\n", addr)

		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Printf("net listen error: %s\n", err.Error())
			return
		}

		if err := s.grpcCli.Serve(lis); err != nil {
			log.Printf("server start error: %s\n", err.Error())
		}
	}()

	// wait system notifiers or cancel func
	<-ctxSignal.Done()
}

func (s *server) Stop(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	s.grpcCli.GracefulStop()

	err := s.router.Shutdown(ctx)
	if err != nil {
		log.Printf("rest api server stop error: %s\n", err.Error())
	}

	if err == nil {
		log.Println("server stopped successfully with no error")
	} else {
		log.Println("server stop done")
	}
}
