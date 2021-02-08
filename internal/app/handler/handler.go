package handler

import (
	"net/http"

	"github.com/Akezhan1/forum/internal/app/service"
)

type Handler struct {
	services *service.Service
}

type route struct {
	Path    string
	Handler http.HandlerFunc
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s}
}

func (h *Handler) InitRouter() *http.ServeMux {
	routes := []route{
		{
			Path:    "/signup",
			Handler: h.SignUp,
		},
		{
			Path:    "/signin",
			Handler: h.SignIn,
		},
		{
			Path:    "/logout",
			Handler: h.LogOut,
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	for _, route := range routes {
		route.Handler = h.CookiesCheckMiddleware(route.Handler)
		mux.HandleFunc(route.Path, route.Handler)
	}

	return mux
}
