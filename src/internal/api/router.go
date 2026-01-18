package api

import (
	"net/http"
)

type CustomRouter struct {
	*http.ServeMux
}

type HTTPHandlerWithErr func(http.ResponseWriter, *http.Request) error

func (r *CustomRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, pattern := r.Handler(req)

	if pattern == "" {
		err := NewHttpError(http.StatusNotFound, "Route not found - "+req.URL.Path)
		HandleHttpError(w, req, err)
		return
	}

	handler.ServeHTTP(w, req)
}

func (r *CustomRouter) HandleFunc(pattern string, handler HTTPHandlerWithErr) {
	r.ServeMux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			HandleHttpError(w, r, err)
		}
	})
}
