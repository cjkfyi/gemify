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

	//
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	//

	// Server initialization
	proxySvr, proxyAddr, err := api.SetupProxy()
	if err != nil {
		log.Fatal(err)
	}
	grpcSvr, grpcLis, grpcAddr, err := api.SetupGRPC()
	if err != nil {
		log.Fatal(err)
	}

	//

	g, gCtx := errgroup.WithContext(ctx)

	errChanProxy := make(chan error)
	errChanGRPC := make(chan error)

	//

	// Launch servers
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

	//

	// Error monitoring
	g.Go(func() error {
		for {
			select {
			case err := <-errChanProxy:
				fmt.Println("Proxy server error:", err)
				cancel() // Call the original cancel function
			case err := <-errChanGRPC:
				fmt.Println("gRPC server error:", err)
				cancel() // Call the original cancel function
			case <-gCtx.Done():
				return gCtx.Err()
			}
		}
	})

	// Shutdown logic
	g.Go(func() error {
		<-gCtx.Done()

		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelShutdown()

		// Collect any errors
		var shutdownErrors []error

		// Shutdown proxy server
		if err := proxySvr.Shutdown(shutdownCtx); err != nil {
			shutdownErrors = append(shutdownErrors, err)
		}

		// Shutdown gRPC server
		grpcSvr.GracefulStop()

		// Close out of datastores
		if err := api.GracefulClosure(); err != nil {
			shutdownErrors = append(shutdownErrors, err)
		} else {
			fmt.Println("Closed out of datastores?!")
		}

		// Return a more appropriate error
		if len(shutdownErrors) > 0 {
			return fmt.Errorf("multiple shutdown errors: %v", shutdownErrors)
		}
		return nil
	})

	// Wait for all goroutines (errgroup)
	if err := g.Wait(); err != nil {
		log.Fatalf("Exit due to err: %v", err)
	}
}
