package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gentlemanautomaton/signaler"
)

func main() {
	// Capture shutdown signals
	shutdown := signaler.New().Capture(os.Interrupt, syscall.SIGTERM)

	// Parse arguments and environment
	c := DefaultConfig
	c.ParseEnv()
	if len(os.Args) > 0 {
		c.ParseArgs(os.Args[1:], flag.ExitOnError)
	}

	// Prepare an http server
	s := &http.Server{
		Addr:    "0.0.0.0:80",
		Handler: http.StripPrefix(c.URLPrefix, http.FileServer(http.Dir(c.DataPath))),
	}

	// Tell the server to stop gracefully when a shutdown signal is received
	stopped := shutdown.Then(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		s.Shutdown(ctx)
	})

	// Always cleanup and wait until the shutdown has completed
	defer stopped.Wait()
	defer shutdown.Trigger()

	// Run the server and print the final result
	fmt.Println(s.ListenAndServe())
}
