package handler

import (
	"html/template"
	"net/http"
)

func (h *Handler) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			tmpl := template.Must(template.ParseFiles("./web/template/index.html"))
			posts, err := h.services.Post.GetAll()
			if err != nil {
				writeResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			tmpl.Execute(w, posts)
		default:
			writeResponse(w, http.StatusBadRequest, "Bad Method")
		}
	}
}
