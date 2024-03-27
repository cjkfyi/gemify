package main

import (
	"log"
	"sync"

	"gemify/api"
)

func main() {

	// Initialize WaitGroup
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("HTTP proxy server panicked: %v", r)
				// recovery for the Proxy server
			}
		}()
		// Start Proxy Server
		if err := api.StartHTTPProxy(); err != nil {
			log.Fatalf("HTTP proxy server exited with error: %v", err)
		}
	}()

	wg.Add(2)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("gRPC server panicked: %v", r)
				// recovery for the gRPC server
			}
		}()
		// Start gRPC Server
		if err := api.StartGRPCServer(); err != nil {
			log.Fatalf("gRPC server exited with error: %v", err)
		}
	}()

	wg.Wait()
}
