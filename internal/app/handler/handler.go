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
			Path:       "/",
			Handler:    h.Index(),
			NeedAuth:   false,
			OnlyUnauth: false,
		},

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

		{
			Path:       "/post/create",
			Handler:    h.CreatePost(),
			NeedAuth:   true,
			OnlyUnauth: false,
		},
		{
			Path:       "/post/rate",
			Handler:    h.RatePost,
			NeedAuth:   true,
			OnlyUnauth: false,
		},
		{
			Path:       "/post/",
			Handler:    h.GetPost(),
			NeedAuth:   false,
			OnlyUnauth: false,
		},
		{
			Path:       "/filter/",
			Handler:    h.Filter(),
			NeedAuth:   false,
			OnlyUnauth: false,
		},

		{
			Path:       "/comment/create",
			Handler:    h.CreateComment,
			NeedAuth:   true,
			OnlyUnauth: false,
		},
		{
			Path:       "/comment/rate",
			Handler:    h.RateComment,
			NeedAuth:   true,
			OnlyUnauth: false,
		},
	}

	mux := http.NewServeMux()

	for _, route := range routes {

		if route.NeedAuth {
			route.Handler = h.NeedAuthMiddleware(route.Handler)
			fmt.Println("Auth", route.Path)
		}

		if route.OnlyUnauth {
			route.Handler = h.OnlyUnauthMiddleware(route.Handler)
			fmt.Println("Unauth", route.Path)
		}

		route.Handler = h.CookiesCheckMiddleware(route.Handler)

		mux.HandleFunc(route.Path, route.Handler)
	}

	fileServer := http.FileServer(http.Dir("./web/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	//mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./web/static"))))
	//mux.Handle("/static/", http.StripPrefix("/post/", http.FileServer(http.Dir("./web/static"))))

	//http.Handle("/", http.FileServer(http.Dir("./web/static")))

	return mux
}
