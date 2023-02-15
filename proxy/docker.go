package proxy

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/docker/docker/client"
)

type Server struct {
	Listener net.Listener
	Client   *client.Client
}

func (s *Server) Serve(ctx context.Context) error {
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

			log.Println("Got request for", r.URL.String())
		},
		Transport: transport,
	}

	router := http.NewServeMux()
	router.Handle("/", proxy)

	router.HandleFunc("/v1.41/images/json", func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusInternalServerError)
	})

	server := http.Server{
		Handler: router,
	}

	defer s.Listener.Close()

	go func() {
		<-ctx.Done()
		timeout, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()
		server.Shutdown(timeout)
	}()

	if err := server.Serve(s.Listener); err != http.ErrServerClosed {
		return err
	}
	log.Println("Shutting down HTTP server")

	return nil
}
