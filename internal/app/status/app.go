package status

import (
	"context"
	_ "embed" // Go 1.16, template embedding
	"fmt"
	"html/template"
	"net/http"

	// "sync"
	"time"

	"go-project-template/internal/pkg/configs"
	"go-project-template/internal/sources"
	"go-project-template/internal/sources/log"

	"github.com/VictoriaMetrics/metrics"
)

//go:embed template/index.html
var tmplIndex string

type handler struct {
	l   log.Logger
	cfg *configs.Root
}

// Run runs the status server: show the current configuration
// and exposes Promethues metrics
func Run(ctx context.Context, cfg *configs.Root, src *sources.Sources, l log.Logger) error {
	h := handler{
		l:   l,
		cfg: cfg,
	}
	defer h.l.Info("exit")

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.indexPage)
	mux.HandleFunc("/health", h.healthCheckPage)
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WritePrometheus(w, true)
	})
	mux.HandleFunc("string", func(w http.ResponseWriter, r *http.Request) {})

	addr := fmt.Sprintf("%v:%v", cfg.HttpServer.Address, cfg.HttpServer.Port)
	h.l.Info("starting on %v", addr)
	s := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			h.l.Error("error to start the HTTP-status server: %v", err)
			return
		}
	}()

	<-ctx.Done()

	h.l.Info("shutting down")
	ctxIn, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := s.Shutdown(ctxIn); err != nil {
		h.l.Error("stoping gracefully: %v", err)
	}

	return nil
}

func (h *handler) indexPage(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			h.l.Error("index page: close body: %v", err)
		}
	}()

	t, err := template.New("index").Parse(tmplIndex)
	if err != nil {
		h.l.Error("parse template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, struct {
		Version     string
		ServiceName string
		Config      *configs.Root
	}{
		Version:     configs.Version,
		ServiceName: configs.ServiceName,
		Config:      h.cfg,
	})

	if err != nil {
		h.l.Error("execute template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) healthCheckPage(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
