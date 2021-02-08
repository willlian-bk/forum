package handler

import (
	"net/http"
)

func (h *Handler) CookiesCheckMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("forum")
		if c != nil {
			if !h.services.IsValidToken(c.Value) {
				c.MaxAge = -1
				http.SetCookie(w, c)
			}
		}
		next.ServeHTTP(w, r)
	})
}
