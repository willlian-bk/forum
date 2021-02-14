package handler

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Akezhan1/forum/internal/app/models"
	uuid "github.com/satori/go.uuid"
)

func (h *Handler) CreatePost() http.HandlerFunc {
	type viewData struct {
		Categories []string
	}

	const maxUploadImage = 20 * 1024 * 1024

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			tmpl := template.Must(template.ParseFiles("./web/template/create_post.html"))
			if categories, err := h.services.Post.GetValidCategories(); err != nil {
				writeResponse(w, http.StatusInternalServerError, err.Error())
			} else {
				tmpl.Execute(w, viewData{categories})
			}
		case "POST":
			c, _ := r.Cookie("forum")
			userID, err := h.services.User.GetUserIDByToken(c.Value)
			if err != nil {
				writeResponse(w, http.StatusForbidden, "Invalid Token")
				return
			}

			r.Body = http.MaxBytesReader(w, r.Body, maxUploadImage)
			if err := r.ParseMultipartForm(maxUploadImage); err != nil {
				writeResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			r.ParseForm()

			post := &models.Post{
				UserID:     userID,
				Title:      r.FormValue("title"),
				Content:    r.FormValue("content"),
				Categories: r.Form["categories"],
			}

			formdata := r.MultipartForm
			files := formdata.File["files"]

			created := false

			for i := range files {
				file, err := files[i].Open()
				if err != nil {
					writeResponse(w, 500, err.Error())
					return
				}

				defer file.Close()
				buff := make([]byte, 512)
				_, err = file.Read(buff)
				if err != nil {
					writeResponse(w, http.StatusInternalServerError, "Something Wrong")
					return
				}
				fileType := http.DetectContentType(buff)
				if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/jpg" && fileType != "image/gif" {
					writeResponse(w, http.StatusInternalServerError, "Invalid File Type")
					return
				}
				_, err = file.Seek(0, io.SeekStart)
				if err != nil {
					writeResponse(w, http.StatusInternalServerError, "Invalid File Type")
					return
				}
				err = os.MkdirAll("./assets/images", os.ModePerm)
				if err != nil {
					writeResponse(w, http.StatusInternalServerError, err.Error())
					return
				}

				imageName := uuid.NewV4().String()
				destImage := fmt.Sprintf("/images/%s%s", imageName, filepath.Ext(files[i].Filename))
				dst, err := os.Create("./assets" + destImage)
				if err != nil {
					writeResponse(w, http.StatusInternalServerError, err.Error())
					return
				}
				defer dst.Close()
				_, err = io.Copy(dst, file)
				if err != nil {
					writeResponse(w, http.StatusInternalServerError, err.Error())
					return
				}

				if !created {
					code, id, err := h.services.Post.Create(post)
					if err != nil {
						writeResponse(w, code, err.Error())
						return
					}

					created = true
					post.ID = id
				}

				if err := h.services.Post.SetImage(post.ID, destImage); err != nil {
					writeResponse(w, http.StatusInternalServerError, err.Error())
					return
				}
			}

			if !created {
				code, id, err := h.services.Post.Create(post)
				if err != nil {
					writeResponse(w, code, err.Error())
					return
				}

				created = true
				post.ID = id
			}

			http.Redirect(w, r, fmt.Sprintf("/post/%d", post.ID), http.StatusFound)
		default:
			writeResponse(w, http.StatusBadRequest, "Bad Method")
		}
	}
}

func (h *Handler) GetPost() http.HandlerFunc {
	type viewData struct {
		Post     *models.Post
		PostID   int
		LoggedIn bool
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			id := getPostIDFromURL(r.URL.Path)
			if post, err := h.services.Post.Get(id); err != nil {
				writeResponse(w, http.StatusBadRequest, err.Error())
			} else {
				tmpl := template.Must(template.ParseFiles("./web/template/view_post.html"))
				ok := IsLoggedUser(r)
				viewData := viewData{post, post.ID, ok}
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

func (h *Handler) Filter() http.HandlerFunc {
	type templateData struct {
		Posts           []*models.Post
		LoggedIn        bool
		ValidCategories []string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			tmpl := template.Must(template.ParseFiles("./web/template/index.html"))
			field := getFiltersFieldFromURL(r.URL.Path)

			userID := 0
			var err error

			c, _ := r.Cookie("forum")
			if c != nil {
				userID, err = h.services.User.GetUserIDByToken(c.Value)
				if err != nil {
					writeResponse(w, http.StatusForbidden, "Invalid Token")
					return
				}
			}

			posts, err := h.services.Post.Filter(field, userID)
			if err != nil {
				if err.Error() == "Unauthorized" {
					http.Redirect(w, r, "/signin", http.StatusFound)
				} else {
					writeResponse(w, http.StatusInternalServerError, err.Error())
				}
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
