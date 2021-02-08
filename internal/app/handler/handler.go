package handler

import (
	"fmt"
	"net/http"

	"github.com/Akezhan1/forum/internal/app/service"
)

type Handler struct {
	services *service.Service
}

type route struct {
	Path       string
	Handler    http.HandlerFunc
	NeedAuth   bool
	OnlyUnauth bool
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s}
}

func (h *Handler) InitRouter() *http.ServeMux {
	routes := []route{
		{
			Path:       "/signup",
			Handler:    h.SignUp,
			NeedAuth:   false,
			OnlyUnauth: true,
		},
		{
			Path:       "/signin",
			Handler:    h.SignIn,
			NeedAuth:   false,
			OnlyUnauth: true,
		},
		{
			Path:       "/logout",
			Handler:    h.LogOut,
			NeedAuth:   true,
			OnlyUnauth: false,
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	for _, route := range routes {
		route.Handler = h.CookiesCheckMiddleware(route.Handler)

		if route.NeedAuth {
			route.Handler = h.NeedAuthMiddleware(route.Handler)
			fmt.Println("Auth", route.Path)
		}

		if route.OnlyUnauth {
			route.Handler = h.OnlyUnauthMiddleware(route.Handler)
			fmt.Println("Unauth", route.Path)
		}

		mux.HandleFunc(route.Path, route.Handler)
	}

	return mux
}
