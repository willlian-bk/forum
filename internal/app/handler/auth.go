package handler

import (
	"html/template"
	"net/http"

	"github.com/Akezhan1/forum/internal/app/models"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
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
			return
		}

		user.ID = id

		writeResponse(w, code, user)
	default:
		writeResponse(w, http.StatusBadRequest, "Bad Method")
	}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.ParseFiles("./web/template/signin.html"))
		tmpl.Execute(w, nil)
	case "POST":
		login := r.FormValue("login")
		password := r.FormValue("password")

		u, err := h.services.User.Authorization(login, password)
		if err != nil {
			writeResponse(w, http.StatusBadRequest, err.Error())
		} else {
			writeResponse(w, http.StatusOK, u)
		}
	default:
		writeResponse(w, http.StatusBadRequest, "Bad Method")
	}
}
