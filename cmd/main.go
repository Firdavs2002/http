package main

import (
	"net"
	"net/http"
	"os"
	"sync"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}

func execute(host string, port string) (err error) {
	mux := http.NewServeMux()
	server := app.NewServer(mux)

	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}

	return srv.ListenAndServe()
}

type handler struct {
	mu       *sync.RWMutex
	handlers map[string]http.HandlerFunc
}

// ServeHTPP обрабатывает все запросы
func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.mu.RLock()
	handler, ok := h.handlers[request.URL.Path]
	h.mu.RUnlock()
	if !ok {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	handler(writer, request)
}
