package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"gemify/api"
)

func main() {

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	//  Init
	proxySvr, proxyAddr, err := api.SetupProxy()
	if err != nil {
		log.Fatal(err)
	}
	grpcSvr, grpcLis, grpcAddr, err := api.SetupGRPC()
	if err != nil {
		log.Fatal(err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	errChanProxy := make(chan error)
	errChanGRPC := make(chan error)

	//  Launch
	g.Go(func() error {
		fmt.Printf("\n  🌱 Proxy Server running @ %v\n", proxyAddr)
		err := proxySvr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errChanProxy <- err
		}
		return nil
	})
	g.Go(func() error {
		fmt.Printf("\n  🌱 gRPC Server running @ %v\n\n", grpcAddr)
		err := grpcSvr.Serve(grpcLis)
		if err != nil {
			errChanGRPC <- err
		}
		return nil
	})

	//  Monitor
	g.Go(func() error {
		for {
			select {
			case err := <-errChanProxy:
				fmt.Println("Proxy server error:", err)
				cancel()
			case err := <-errChanGRPC:
				fmt.Println("gRPC server error:", err)
				cancel()
			case <-gCtx.Done():
				return gCtx.Err()
			}
		}
	})

	//  Shutdown
	g.Go(func() error {
		<-gCtx.Done()

		shutdownCtx, cancelShutdown := context.WithTimeout(
			context.Background(),
			30*time.Second,
		)
		defer cancelShutdown()

		var shutdownErrors []error

		if err := proxySvr.Shutdown(shutdownCtx); err != nil {
			shutdownErrors = append(shutdownErrors, err)
		}

		grpcSvr.GracefulStop()

		// if err := store.GracefulClosure(); err != nil {
		// 	shutdownErrors = append(shutdownErrors, err)
		// } else {
		// 	fmt.Println(" Closed out of datastores!")
		// }

		if len(shutdownErrors) > 0 {
			return fmt.Errorf("multiple shutdown errors: %v", shutdownErrors)
		}
		return nil
	})

	//  Wait...
	if err := g.Wait(); err != nil {
		log.Fatalf("Exit due to err: %v", err)
	}
}
