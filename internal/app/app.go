package app

import (
	"context"
	"fmt"
	"github.com/stawwkom/auth_service/internal/interceptor"
	"github.com/stawwkom/auth_service/internal/logger"
	descAccess "github.com/stawwkom/auth_service/pkg/access_v1"
	descAuth "github.com/stawwkom/auth_service/pkg/auth_login"
	desc "github.com/stawwkom/auth_service/pkg/auth_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"
	"net"
	"net/http"
	"os"

	"crypto/tls"
	runtimes "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stawwkom/auth_service/internal/config"
	"github.com/stawwkom/auth_service/internal/repository"
	repo "github.com/stawwkom/auth_service/internal/repository/auth"
	serv "github.com/stawwkom/auth_service/internal/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	defaultCertPath = "certs/service.pem"
	defaultKeyPath  = "certs/service.key"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func getCertPaths() (string, string) {
	certPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("KEY_PATH")

	if certPath == "" {
		// Проверяем, запущено ли приложение из cmd/auth_server
		if _, err := os.Stat("../../certs/service.pem"); err == nil {
			certPath = "../../certs/service.pem"
			keyPath = "../../certs/service.key"
		} else {
			certPath = defaultCertPath
			keyPath = defaultKeyPath
		}
	}

	return certPath, keyPath
}

func (a *App) Run(ctx context.Context) error {
	logger.Info("🚀 gRPC сервер запущен",
		zap.String("address", a.serviceProvider.GRPCAddr()),
		zap.String("log_level", config.Cfg.Log.Level),
	)

	listener, err := net.Listen("tcp", a.serviceProvider.GRPCAddr())
	if err != nil {
		return err
	}

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		if err := a.grpcServer.Serve(listener); err != nil {
			logger.Error("❌ gRPC сервер остановлен", zap.Error(err))
		}
	}()

	// Запускаем HTTP сервер в отдельной горутине
	go func() {
		certPath, keyPath := getCertPaths()
		logger.Info("🌐 HTTP сервер запущен",
			zap.String("address", a.serviceProvider.HTTPAddr()),
		)
		if err := a.httpServer.ListenAndServeTLS(certPath, keyPath); err != nil && err != http.ErrServerClosed {
			logger.Error("❌ HTTP сервер остановлен", zap.Error(err))
		}
	}()

	// Ожидаем завершения контекста
	<-ctx.Done()

	// Завершаем gRPC сервер
	logger.Info("⏹ Останавливаем gRPC сервер...")
	a.grpcServer.GracefulStop()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	steps := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
	}

	for _, step := range steps {
		if err := step(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	return config.Load()
}

func (a *App) initServiceProvider(ctx context.Context) error {
	db, err := repository.NewPostgresDB(ctx, config.Cfg)
	if err != nil {
		return err
	}

	authRepo := repo.NewRepository(db)
	authService := serv.NewAuthService(authRepo)

	a.serviceProvider = &serviceProvider{
		authService: authService,
		dbClose: func() error {
			db.Close()
			return nil
		},
	}
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	certPath, keyPath := getCertPaths()
	creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		return fmt.Errorf("failed to load TLS cert: %v", err)
	}

	a.grpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(
			interceptor.LogInterceptor,
			interceptor.ValidateInterceptor,
		),
	)
	fmt.Println("Validate running")
	reflection.Register(a.grpcServer)
	//desc.RegisterUserAPIServer(a.grpcServer, api.NewServer(a.serviceProvider.authService))
	desc.RegisterUserAPIServer(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	descAccess.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessI(ctx))
	descAuth.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthI(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	certPath, _ := getCertPaths()
	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		return fmt.Errorf("failed to load TLS cert: %v", err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	mux := runtimes.NewServeMux()

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	a.httpServer = &http.Server{
		Addr:      a.serviceProvider.HTTPAddr(),
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	if err := desc.RegisterUserAPIHandlerFromEndpoint(
		ctx, mux, a.serviceProvider.GRPCAddr(), opts,
	); err != nil {
		return err
	}

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPAddr(),
		Handler: mux,
	}

	return nil
}

func (a *App) Close() {
	logger.Warn("⏹ Закрытие приложения...")
	if err := a.serviceProvider.dbClose(); err != nil {
		logger.Error("Ошибка при закрытии БД", zap.Error(err))
	}
	if a.grpcServer != nil {
		logger.Info("⏹ gRPC остановка...")
		a.grpcServer.GracefulStop()
	}
	if a.httpServer != nil {
		logger.Info("⏹ HTTP остановка...")
		if err := a.httpServer.Shutdown(context.Background()); err != nil {
			logger.Error("Ошибка при остановке HTTP сервера", zap.Error(err))
		}
	}
}
