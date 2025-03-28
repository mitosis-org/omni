package app

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/omni-network/omni/lib/contracts"
	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/log"
	"github.com/omni-network/omni/solver/types"

	"github.com/rs/cors"
)

type (
	JSONResponse      = types.JSONResponse
	JSONErrorResponse = types.JSONErrorResponse
)

// removeBUG removes [BUG] from the error messages, so they are not included in responses to users.
func removeBUG(s string) string { return strings.ReplaceAll(s, "[BUG]", "") }

func writeJSON(ctx context.Context, w http.ResponseWriter, res JSONResponse) {
	w.WriteHeader(res.StatusCode())
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error(ctx, "[BUG] error writing /quote response", err)
	}
}

// serveAPI starts the API server, returning a async error.
func serveAPI(address string, endpoints map[string]http.Handler) <-chan error {
	errChan := make(chan error)
	go func() {
		mux := http.NewServeMux()

		endpoints["/live"] = newLiveHandler()
		endpoints["/"] = newLiveHandler() // Also serve live from root for easy health checks

		for endpoint, handler := range endpoints {
			mux.Handle(endpoint, instrumentHandler(endpoint, handler))
		}

		c := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Origin", "Content-Type", "Accept"},
		})

		srv := &http.Server{
			Addr:              address,
			ReadHeaderTimeout: 5 * time.Second,
			IdleTimeout:       5 * time.Second,
			WriteTimeout:      5 * time.Second,
			Handler:           c.Handler(mux),
		}
		errChan <- errors.Wrap(srv.ListenAndServe(), "serve api")
	}()

	return errChan
}

// newContractsHandler returns a http handler that returns the contract address for `network`.
func newContractsHandler(addrs contracts.Addresses) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rr *http.Request) {
		ctx := rr.Context()

		w.Header().Set("Content-Type", "application/json")

		writeJSON(ctx, w, types.ContractsResponse{
			Portal:    addrs.Portal.Hex(),
			Inbox:     addrs.SolverNetInbox.Hex(),
			Outbox:    addrs.SolverNetOutbox.Hex(),
			Middleman: addrs.SolverNetMiddleman.Hex(),
		})
	})
}

// newLiveHandler returns a http handler that always return 200s.
// It indicates the API server is "live" (up and running).
func newLiveHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

// instrumentHandler wraps an http.Handler, instrumenting the latency and error rate.
func instrumentHandler(endpoint string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		iw := &instrumentWriter{ResponseWriter: w}

		handler.ServeHTTP(iw, r)

		apiLatency.WithLabelValues(endpoint).Observe(time.Since(t0).Seconds())
		apiErrors.WithLabelValues(endpoint).Add(iw.ErrorCount())
	})
}

// instrumentWriter wraps the response writer, tracking the status code written.
type instrumentWriter struct {
	http.ResponseWriter
	status int
}

func (w *instrumentWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *instrumentWriter) ErrorCount() float64 {
	if w.status < 400 {
		return 0 // 200 or unset is not an error
	}

	return 1
}
