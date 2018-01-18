package main

import (
	"fmt"
	"net/http"
)

type Handler struct{}

// http(s) handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// log all incoming request
	ch <- fmt.Sprintf(
		"%s - %s (%s) %s", r.RemoteAddr, r.Method, r.Proto, r.RequestURI,
	)
	// getting configured hosts
	handler, ok := domains[r.Host]
	if !ok {
		http.NotFound(w, r)
		return
	}

	if useTLS && r.TLS == nil {
		// redirect to https port
		http.Redirect(
			w, r, fmt.Sprintf("https://%s/%s", r.Host, r.RequestURI), 301,
		)
		return
	}

	if len(r.RequestURI) == 1 {
		r.RequestURI = "/index.html"
	}

	// handling the request
	handler.ServeHTTP(w, r)
}
