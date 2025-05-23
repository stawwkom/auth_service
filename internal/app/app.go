package app

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stawwkom/auth_service/internal/interceptor"
	"github.com/stawwkom/auth_service/internal/logger"
	"github.com/stawwkom/auth_service/internal/metric"
	descAccess "github.com/stawwkom/auth_service/pkg/access_v1"
	descAuth "github.com/stawwkom/auth_service/pkg/auth_login"
	desc "github.com/stawwkom/auth_service/pkg/auth_v1"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials"
	"log"
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
	logger.Info("🚀 gRPC сервер запущен")

	listener, err := net.Listen("tcp", a.serviceProvider.GRPCAddr())
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return a.grpcServer.Serve(listener)
	})

	g.Go(func() error {
		certPath, keyPath := getCertPaths()
		return a.httpServer.ListenAndServeTLS(certPath, keyPath)
	})

	g.Go(func() error {
		return runPrometheus(ctx) // передать ctx
	})

	// Ждём завершения или ошибки
	if err := g.Wait(); err != nil {
		logger.Error("⛔ Ошибка в одной из горутин", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	steps := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initMetrics,
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
			interceptor.MetricsInterceptor,
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

func (a *App) initMetrics(ctx context.Context) error {
	return metric.Init(ctx)
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

func runPrometheus(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    "0.0.0.0:2112",
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		_ = server.Shutdown(context.Background())
	}()

	log.Println("Prometheus server is running on 0.0.0.0:2112")
	return server.ListenAndServe()
}
