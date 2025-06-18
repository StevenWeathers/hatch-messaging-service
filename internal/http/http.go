package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

func (o *Endpoints) handle() {
	// k8s healthcheck /healthz as per convention
	o.router.HandleFunc("GET /healthz", o.handleHealthz)
}

func (o *Endpoints) ListenAndServe(ctx context.Context) error {
	slog.InfoContext(ctx, fmt.Sprintf("Listening on %s", o.config.ListenAddress))
	return http.ListenAndServe(o.config.ListenAddress, o.router)
}

func New(config Config) *Endpoints {
	router := http.NewServeMux()
	e := &Endpoints{
		config: config,
		router: router,
	}
	e.handle()
	return e
}
