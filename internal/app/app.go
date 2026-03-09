// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	authGrpc "github.com/KimNattanan/go-chat-backend/internal/auth/handler/grpc"
	authRest "github.com/KimNattanan/go-chat-backend/internal/auth/handler/rest"
	authPb "github.com/KimNattanan/go-chat-backend/internal/auth/proto/v1"
	authPersistent "github.com/KimNattanan/go-chat-backend/internal/auth/repo/persistent"
	authUseCase "github.com/KimNattanan/go-chat-backend/internal/auth/usecase/auth"

	profileAmqpRpc "github.com/KimNattanan/go-chat-backend/internal/profile/handler/amqp_rpc"
	profileGrpc "github.com/KimNattanan/go-chat-backend/internal/profile/handler/grpc"
	profileRest "github.com/KimNattanan/go-chat-backend/internal/profile/handler/rest"
	profilePb "github.com/KimNattanan/go-chat-backend/internal/profile/proto/v1"
	profilePersistent "github.com/KimNattanan/go-chat-backend/internal/profile/repo/persistent"
	profileUseCase "github.com/KimNattanan/go-chat-backend/internal/profile/usecase/profile"

	chatRest "github.com/KimNattanan/go-chat-backend/internal/chat/handler/rest"
	chatPersistent "github.com/KimNattanan/go-chat-backend/internal/chat/repo/persistent"
	membershipUseCase "github.com/KimNattanan/go-chat-backend/internal/chat/usecase/membership"
	roomUseCase "github.com/KimNattanan/go-chat-backend/internal/chat/usecase/room"

	messageRest "github.com/KimNattanan/go-chat-backend/internal/message/handler/rest"
	messagePersistent "github.com/KimNattanan/go-chat-backend/internal/message/repo/persistent"
	messageUseCase "github.com/KimNattanan/go-chat-backend/internal/message/usecase/message"

	"github.com/KimNattanan/go-chat-backend/internal/platform/config"
	"github.com/KimNattanan/go-chat-backend/internal/platform/middleware"

	"github.com/KimNattanan/go-chat-backend/pkg/grpcserver"
	"github.com/KimNattanan/go-chat-backend/pkg/httpserver"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/postgres"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
	"github.com/KimNattanan/go-chat-backend/pkg/redisclient"
	"github.com/KimNattanan/go-chat-backend/pkg/token"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	jwtMaker := token.NewJWTMaker(cfg.JWT.Secret)

	// Repository
	pg, err := postgres.New(cfg.DB.DSN)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	rdb := redisclient.New(cfg.Redis.Address, cfg.Redis.Password, cfg.Redis.DB)

	// gRPC Client
	grpcClientConn, err := grpc.NewClient("localhost:"+cfg.GRPC.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - grpc.NewClient: %w", err))
	}
	defer grpcClientConn.Close()

	authGrpcClient := authPb.NewAuthServiceClient(grpcClientConn)
	profileGrpcClient := profilePb.NewProfileServiceClient(grpcClientConn)

	// RabbitMQ conn
	rmqClient := rabbitmq.New(cfg.RMQ.URL, "queue1")
	if err := rmqClient.Connect(); err != nil {
		l.Fatal(fmt.Errorf("app - Run - rmqClient.New: %w", err))
	}

	// Use-Case
	authUseCase := authUseCase.New(
		authPersistent.NewUserRepo(pg.DB),
		authPersistent.NewSessionRepo(rdb),
		profileGrpcClient,
		rmqClient,
		jwtMaker,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)

	profileUseCase := profileUseCase.New(
		profilePersistent.NewProfileRepo(pg.DB),
	)

	roomUseCase := roomUseCase.New(
		chatPersistent.NewRoomRepo(pg.DB),
		rmqClient,
	)
	membershipUseCase := membershipUseCase.New(
		chatPersistent.NewMembershipRepo(pg.DB),
		authGrpcClient,
	)

	messageUseCase := messageUseCase.New(
		messagePersistent.NewMessageRepo(pg.DB),
	)

	// RabbitMQ Consumer
	rmqRouter := make(map[string]rabbitmq.Handler)
	profileAmqpRpc.NewRouter(rmqRouter, profileUseCase, l)

	// gRPC Server
	grpcServer := grpcserver.New(l, grpcserver.Port(cfg.GRPC.Port))
	authGrpc.NewRouter(grpcServer.App, authUseCase, l)
	profileGrpc.NewRouter(grpcServer.App, profileUseCase, l)
	reflection.Register(grpcServer.App)

	// Middleware
	jwtMiddleware := middleware.JWTMiddleware(l, cfg, jwtMaker, authGrpcClient)

	// HTTP Server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))
	httpServer.Echo.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Accept", "Content-Type", "Origin", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))
	httpServer.Echo.Use(middleware.Logger(l))
	httpServer.Echo.Use(middleware.Recovery(l))
	authRest.NewRouter(httpServer.Echo, cfg, authUseCase, l, jwtMiddleware)
	profileRest.NewRouter(httpServer.Echo, cfg, profileUseCase, l, jwtMiddleware)
	chatRest.NewRouter(httpServer.Echo, cfg, roomUseCase, membershipUseCase, l, jwtMiddleware)
	messageRest.NewRouter(httpServer.Echo, cfg, messageUseCase, l, jwtMiddleware)

	// Start servers
	if err := rmqClient.Consume(2, rmqRouter); err != nil {
		l.Error("app - Run - rmqClient.Consume: %w", err)
	}
	grpcServer.Start()
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-grpcServer.Notify():
		l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = grpcServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}

	err = rmqClient.Close()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - rmqClient.Close: %w", err))
	}
}
