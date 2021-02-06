package handler

import "net/http"

type Handler struct {
	//services *service.Service
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) InitRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	return mux
}
