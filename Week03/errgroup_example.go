package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
)


var addr = "0.0.0.0:8080"

type Handel struct {
}

func (c *Handel) ServeHTTP(w http.ResponseWriter, req *http.Request) {
}

func (c Handel) Ping(w http.ResponseWriter, req *http.Request) {
}

func serveHttp(ctx context.Context) error {
	defer func() {
		fmt.Println("http server prepare exit")
	}()

	srv := &http.Server{
		Addr:    addr,
		Handler: &Handel{},
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(ctx)
	}()

	err := srv.ListenAndServe()

	return fmt.Errorf("serve http failed: %w", err)
}

func handleSignal(ctx context.Context) error {
	defer func() {
		fmt.Println("handle signal prepare exit")
	}()

	signC := make(chan os.Signal)
	signal.Notify(signC, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sign := <-signC:
		err := fmt.Errorf("handle signal: %v", sign)
		return err

	case <-ctx.Done():
		return nil
	}
}

func main() {

	
	group, ctx := errgroup.WithContext(context.Background())
	group.Go(func() error {
		return serveHttp(ctx)
	})

	group.Go(func() error {
		return handleSignal(ctx)
	})

	fmt.Println("server run at: ", addr)
	err := group.Wait()
	fmt.Println("sever recv error", err)
}