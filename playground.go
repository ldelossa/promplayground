package promplayground

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// package globals shared between client and server
const (
	ListenAddr = "localhost:8089"
	DBConnStr  = "host=prom-db port=5432 user=prom dbname=prom sslmode=disable"
	HTTPStatusArray = [...]int{
		http.StatusOK,
		http.StatusInternalServerError,
		http.StatusNotFound,
		http.StatusFound,
		http.StatusCreated,
		http.Status
	}
)

// guages utilized in handlers
var (
	histogram = promauto.NewGaugeVec(
		prometheus.HistogramOpts{
			Name:    "http_response_latency",
			Namespace: "promgplayground",
			Help:    "A histograph of http response latency distributions",
			Buckets: []float64{.030, .1, .5, 1, 10, 30, 60},
		},
	)
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests",
			Namespace: "promgplayground",
			Help: "A monotonically increasing count of requests with useful tags"
		},
		[]string{"path", ""}	
	)
)

// The RandoHandler returns an http handler that is unpredictable.
// The response status, response size, and latency to respond are all
// randomized.
//
// Writing an small tool which smacks this endpoint in a loop will generate 
// generate a random distribution of metric points to play with. 
func RandoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Unix()
		
	}
}
