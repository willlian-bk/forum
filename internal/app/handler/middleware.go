package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) CookiesCheckMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("forum")
		fmt.Println("cookies")
		if c != nil {
			if !h.services.IsValidToken(c.Value) {
				fmt.Println("cookies check")
				c.MaxAge = -1
				http.SetCookie(w, c)
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) OnlyUnauthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("forum")
		fmt.Println("unauth")
		if c != nil {
			http.Redirect(w, r, "/", 301)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (h *Handler) NeedAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("forum")
		fmt.Println("auth")
		if c == nil {
			http.Redirect(w, r, "/signin", 301)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
