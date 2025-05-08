package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/app"
	pb "homework9/internal/ports/grpc"
	ser "homework9/internal/ports/grpc/service"
	"homework9/internal/ports/httpgin"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	root, cancel := context.WithCancel(context.Background())
	er, ctx := errgroup.WithContext(root)

	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	er.Go(func() error {
		select {
		case <-sigQuit:
			log.Println("captured shutdown signal")
			cancel()
			return fmt.Errorf("interrupt signal received")
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	ap := app.NewApp(adrepo.New())
	er.Go(func() error {
		lis, err := net.Listen("tcp", ":8080")
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
		httpServer := httpgin.NewHTTPServer(":1010", ap)

		errCh := make(chan error, 1)
		defer func() {
			ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel2()
			if err := httpServer.Shutdown(ctx2); err != nil {
				log.Printf("failed to shutdown http server: %v", err)
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
		log.Printf("shutting down services:%e \n", err)
	}
	log.Println("goodbye")
}
