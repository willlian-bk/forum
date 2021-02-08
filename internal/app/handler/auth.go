package handler

import (
	"html/template"
	"net/http"

	"github.com/Akezhan1/forum/internal/app/models"
)

func (h *Handler) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			tmpl := template.Must(template.ParseFiles("./web/template/signup.html"))
			tmpl.Execute(w, nil)
		case "POST":
			user := &models.User{
				Email:    r.FormValue("email"),
				Username: r.FormValue("username"),
				Password: r.FormValue("password"),
				Role:     "user",
			}

			code, id, err := h.services.User.Create(user)
			if err != nil {
				writeResponse(w, code, err.Error())
				http.Redirect(w, r, "/signup", 302)
				return
			}

			user.ID = id

			writeResponse(w, code, user)
		default:
			writeResponse(w, http.StatusBadRequest, "Bad Method")
		}
	}
}

func (h *Handler) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			tmpl := template.Must(template.ParseFiles("./web/template/signin.html"))
			tmpl.Execute(w, nil)
		case "POST":
			login := r.FormValue("login")
			password := r.FormValue("password")

			session, err := h.services.User.Authorization(login, password)
			if err != nil {
				writeResponse(w, http.StatusBadRequest, err.Error())
			} else {
				http.SetCookie(w, &http.Cookie{
					Name:    "forum",
					Path:    "/",
					Value:   session.Token,
					Expires: session.ExpTime,
				})
				writeResponse(w, http.StatusOK, "OK")
			}
		default:
			writeResponse(w, http.StatusBadRequest, "Bad Method")
		}
	}
}

func (h *Handler) LogOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			cookie, err := r.Cookie("forum")
			if err != nil {
				writeResponse(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			if err := h.services.User.Logout(cookie.Value); err != nil {
				writeResponse(w, http.StatusInternalServerError, err.Error())
			} else {
				cookie.MaxAge = -1
				http.SetCookie(w, cookie)
				writeResponse(w, http.StatusOK, "OK")
			}
		default:
			writeResponse(w, http.StatusBadRequest, "Bad Method")
		}
	}
}
