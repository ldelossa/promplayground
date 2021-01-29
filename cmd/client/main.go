package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	play "github.com/ldelossa/promplayground"
	"golang.org/x/sync/semaphore"
)

var (
	HTTPMethodArray = [...]string{
		http.MethodGet,
		http.MethodDelete,
		http.MethodPut,
		http.MethodPost,
	}
)

func main() {
	cr := flag.Int("cr", 2, "the number of concurrent requests to make")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		log.Printf("interrupt received, gracefully shutting down")
		cancel()
	}()

	sem := semaphore.NewWeighted(int64(*cr))
	for {
		sem.Acquire(ctx, 1)
		go func() {
			defer sem.Release(1)
			r := &http.Request{
				URL: &url.URL{
					Scheme: "http",
					Host:   play.ListenAddr,
					Path:   play.RandoPath,
				},
				Method: HTTPMethodArray[rand.Int()%len(HTTPMethodArray)],
			}
			log.Printf("making request %s", r.URL)
			_, err := http.DefaultClient.Do(r)
			if err != nil {
				log.Fatal(err)
			}
		}()

		if ctx.Err() != nil {
			os.Exit(0)
		}
	}
}
