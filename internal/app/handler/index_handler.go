package handler

import (
	"html/template"
	"net/http"

	"github.com/Akezhan1/forum/internal/app/models"
)

func (h *Handler) Index() http.HandlerFunc {
	type templateData struct {
		Posts           []*models.Post
		LoggedIn        bool
		ValidCategories []string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if r.URL.Path != "/" {
				writeResponse(w, http.StatusNotFound, "Page Not Found")
				return
			}
			tmpl := template.Must(template.ParseFiles("./web/template/index.html"))
			posts, err := h.services.Post.GetAll()
			if err != nil {
				writeResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			validcategories, err := h.services.Post.GetValidCategories()
			if err != nil {
				writeResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			ok := IsLoggedUser(r)

			tmpl.Execute(w, templateData{posts, ok, validcategories})
		default:
			writeResponse(w, http.StatusBadRequest, "Bad Method")
		}
	}
}
