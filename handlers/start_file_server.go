package handlers

import (
	"net/http"
	"os"
)

func StartFileServer(w http.ResponseWriter, r *http.Request) {
	webDir := "web"
	if _, err := os.Stat(webDir + r.RequestURI); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		http.FileServer(http.Dir(webDir)).ServeHTTP(w, r)
	}
}
