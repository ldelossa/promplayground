package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	play "github.com/ldelossa/promplayground"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var Mux http.ServeMux
	Mux.Handle("/metrics", promhttp.Handler())
	Mux.Handle("/rando", play.RandoHandler())

	var MetricsServer = http.Server{
		Addr:    play.ListenAddr,
		Handler: &Mux,
	}

	log.Printf("launching metrics server on %s%s", play.ListenAddr, "/metrics")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := MetricsServer.ListenAndServe()
		switch err {
		case http.ErrServerClosed:
		case nil:
		default:
			log.Printf("error while launching metrics server: %v", err)
			cancel()
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	select {
	case <-sigint:
		log.Printf("interrupt received, gracefully shutting down")
		tctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		if err := MetricsServer.Shutdown(tctx); err != nil {
			log.Fatalf("failed to gracefully shutdown: %v", err)
		}
		os.Exit(0)
	case <-ctx.Done():
		os.Exit(1)
	}
}
