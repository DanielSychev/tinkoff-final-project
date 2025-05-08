package main

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/app"
	"homework9/internal/config"
	pb "homework9/internal/ports/grpc"
	ser "homework9/internal/ports/grpc/service"
	"homework9/internal/ports/httpgin"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	//logger2 := logger.Sugar()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal("config fail", zap.Error(err))
	}

	root, cancel := context.WithCancel(context.Background())
	er, ctx := errgroup.WithContext(root)
	ctx = context.WithValue(ctx, "logger", logger)

	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	er.Go(func() error {
		select {
		case <-sigQuit:
			logger.Info("captured shutdown signal")
			cancel()
			return fmt.Errorf("interrupt signal received")
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	ap := app.NewApp(adrepo.New())
	er.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort))
		if err != nil {
			return fmt.Errorf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterAdServiceServer(grpcServer, ser.NewMyServer(ap))

		errCh := make(chan error, 1)
		defer func() {
			grpcServer.GracefulStop()
			_ = lis.Close()
			close(errCh)
		}()

		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				errCh <- fmt.Errorf("failed to serve gprc server: %v", err)
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return err
		}
	})

	er.Go(func() error {
		httpServer := httpgin.NewHTTPServer(ctx, fmt.Sprintf(":%d", cfg.RestPort), ap)

		errCh := make(chan error, 1)
		defer func() {
			ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel2()
			if err := httpServer.Shutdown(ctx2); err != nil {
				logger.Error("failed to shutdown http server", zap.Error(err))
			}
			close(errCh)
		}()
		go func() {
			if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				errCh <- fmt.Errorf("failed to serve http server: %v", err)
			}
		}()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return err
		}
	})

	if err := er.Wait(); err != nil {
		logger.Error("shutting down services", zap.Error(err))
	}
	logger.Info("goodbye")
}
