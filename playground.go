package promplayground

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// package globals shared between client and server
const (
	ListenAddr = "0.0.0.0:8089"
	DBConnStr  = "host=prom-db port=5432 user=prom dbname=prom sslmode=disable"
	RandoPath  = "/rando"
)

var (
	HTTPStatusArray = [...]string{
		strconv.Itoa(http.StatusOK),
		strconv.Itoa(http.StatusInternalServerError),
		strconv.Itoa(http.StatusNotFound),
		strconv.Itoa(http.StatusCreated),
	}
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      "http_response_duration_seconds",
			Namespace: "promgplayground",
			Help:      "A histograph of http response duration distribution",
			Buckets:   []float64{.030, .1, .5, 1, 10, 30, 60},
		},
		[]string{"path", "method", "status"},
	)
	RequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      "http_response_size_bytes",
			Namespace: "promgplayground",
			Help:      "A histograph of http response size distribution",
			Buckets:   []float64{500, 1000, 350000, 500000, 1000000},
		},
		[]string{"path", "method", "status"},
	)
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "http_requests_total",
			Namespace: "promgplayground",
			Help:      "A monotonically increasing count of requests with useful tags",
		},
		[]string{"path", "method", "status"},
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
		randDuration := rand.Int() % 120 // seconds
		randResponseSize := rand.Int()   // bytes
		randStatus := rand.Int() % len(HTTPStatusArray)
		time.Sleep(time.Duration(randDuration) * time.Second)

		RequestSize.WithLabelValues(RandoPath, r.Method, HTTPStatusArray[randStatus]).Observe(float64(randResponseSize))
		RequestDuration.WithLabelValues(RandoPath, r.Method, HTTPStatusArray[randStatus]).Observe(float64(randDuration))
		RequestCounter.WithLabelValues(RandoPath, r.Method, HTTPStatusArray[randStatus]).Add(1)

		status, _ := strconv.Atoi(HTTPStatusArray[randStatus])
		w.WriteHeader(status)
	}
}
