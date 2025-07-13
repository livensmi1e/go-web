package platform

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunMonitor(cfg *Config) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		//nolint:errcheck
		fmt.Fprintln(w, "OK")
	})
	mux.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:         cfg.MonitorServerAddr(),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return server.ListenAndServe()
}
