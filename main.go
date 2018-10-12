package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gorilla/handlers"

	"github.com/gentlemanautomaton/signaler"
)

func main() {
	// Print start and stop messages
	fmt.Printf("Starting web server.\n")
	defer fmt.Printf("Stopped web server.\n")

	// Capture shutdown signals
	shutdown := signaler.New().Capture(os.Interrupt, syscall.SIGTERM)

	// Parse arguments and environment
	c := DefaultConfig
	c.ParseEnv()
	if len(os.Args) > 0 {
		c.ParseArgs(os.Args[1:], flag.ExitOnError)
	}

	// Prepare the file system
	fs := http.FileSystem(http.Dir(c.DataPath))
	if !c.DirectoryListing {
		fs = filesOnlyFilesystem{fs}
	}

	// Prepare the handler
	handler := http.FileServer(fs)
	handler = http.StripPrefix(c.URLPrefix, handler)
	if c.Logging {
		handler = handlers.LoggingHandler(os.Stdout, handler)
	}

	// Prepare an http server
	s := &http.Server{
		Addr:    "0.0.0.0:80",
		Handler: handler,
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

type filesOnlyFilesystem struct {
	fs http.FileSystem
}

func (fs filesOnlyFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	if stat, _ := f.Stat(); stat.IsDir() {
		return nil, os.ErrPermission
	}
	return f, nil
}
