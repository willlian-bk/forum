package handler

import (
	"fmt"
	"net/http"
	"time"
)

func (h *Handler) CookiesCheckMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("forum")
		fmt.Println("cookie")
		if c != nil {
			fmt.Println("cookie present")
			if !h.services.IsValidToken(c.Value) {
				fmt.Println("cookie not valid")
				c.MaxAge = -1
				c.Expires = time.Unix(0, 0)
				http.SetCookie(w, c)
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) OnlyUnauthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("forum")
		if c != nil {
			http.Redirect(w, r, "/", 302)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (h *Handler) NeedAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("forum")
		if c == nil {
			http.Redirect(w, r, "/signin", 302)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
