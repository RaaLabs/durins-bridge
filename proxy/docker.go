package proxy

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/docker/docker/client"
	"golang.org/x/exp/slog"
)

type Server struct {
	Listener net.Listener
	Client   *client.Client
}

func (s *Server) Serve(ctx context.Context) error {
	server := http.Server{
		Handler: s.createProxyHandler(),
	}

	defer s.Listener.Close()

	go func() {
		<-ctx.Done()
		timeout, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()
		slog.Info("Shutting down HTTP server")
		server.Shutdown(timeout)
	}()

	if err := server.Serve(s.Listener); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) createProxyHandler() http.Handler {
	dialer := s.Client.Dialer()

	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return dialer(ctx)
		},
	}

	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = "docker.sock"
		},
		Transport: transport,
	}

	return &OverridesHandler{
		Proxy: proxy,
		Overrides: []OverrideHandler{
			&PullOverride{
				Client: s.Client,
			},
		},
	}
}
