package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func writeResponse(w http.ResponseWriter, code int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(code)
	data = append(data, '\n')
	w.Write(data)
}

func getPostIDFromURL(url string) int {
	idStr := strings.TrimPrefix(url, "/post/")

	if idStr == "" {
		return -1
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1
	}

	return id
}

func getFiltersFieldFromURL(url string) string {
	return strings.Title(strings.TrimPrefix(url, "/filter/"))
}
