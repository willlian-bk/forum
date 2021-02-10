package handler

import (
	"html/template"
	"net/http"

	"github.com/Akezhan1/forum/internal/app/models"
)

func (h *Handler) CreatePost() http.HandlerFunc {
	type viewData struct {
		Categories []string
	}

	var data viewData

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			tmpl := template.Must(template.ParseFiles("./web/template/create_post.html"))
			if categories, err := h.services.Post.GetValidCategories(); err != nil {
				writeResponse(w, http.StatusInternalServerError, err.Error())
			} else {
				data.Categories = categories
				tmpl.Execute(w, data)
			}
		case "POST":
			c, _ := r.Cookie("forum")
			userID, err := h.services.User.GetUserIDByToken(c.Value)
			if err != nil {
				writeResponse(w, http.StatusForbidden, "Invalid Token")
				return
			}

			r.ParseForm()
			post := &models.Post{
				UserID:     userID,
				Title:      r.FormValue("title"),
				Content:    r.FormValue("content"),
				Categories: r.Form["categories"],
			}

			code, id, err := h.services.Post.Create(post)
			if err != nil {
				writeResponse(w, code, err.Error())
			} else {
				post.ID = id
				writeResponse(w, code, post)
			}
		default:
			writeResponse(w, http.StatusBadRequest, "Bad Method")
		}
	}
}

func (h *Handler) GetPost() http.HandlerFunc {
	type viewData struct {
		Post   *models.Post
		PostID int
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			id := getPostIDFromURL(r.URL.Path)
			if post, err := h.services.Post.Get(id); err != nil {
				writeResponse(w, http.StatusBadRequest, err.Error())
			} else {
				tmpl := template.Must(template.ParseFiles("./web/template/view_post.html"))
				viewData := viewData{post, post.ID}
				tmpl.Execute(w, viewData)
			}
		default:
			writeResponse(w, http.StatusBadRequest, "Bad Method")
		}
	}
}

func (h *Handler) RatePost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		postID := r.FormValue("post_id")
		types := r.FormValue("type")

		c, _ := r.Cookie("forum")
		userID, err := h.services.User.GetUserIDByToken(c.Value)
		if err != nil {
			writeResponse(w, http.StatusForbidden, "Invalid Token")
			return
		}

		if err := h.services.Post.EstimatePost(postID, userID, types); err != nil {
			writeResponse(w, http.StatusBadRequest, err.Error())
		} else {
			http.Redirect(w, r, "/post/"+postID, http.StatusFound)
		}
	default:
		writeResponse(w, http.StatusBadRequest, "Bad Method")
	}
}
