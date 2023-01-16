package api

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func (api *API) ImgProfileView(w http.ResponseWriter, r *http.Request) {
	// View with response image `img-avatar.png` from path `assets/images`

	convFileBytes, err := ioutil.ReadFile("./assets/images/img-avatar.png")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("File not found"))
		return
	}

	w.Write(convFileBytes)
}

func (api *API) ImgProfileUpdate(w http.ResponseWriter, r *http.Request) {
	// Update image `img-avatar.png` from path `assets/images`

	// r.ParseForm()
	r.ParseMultipartForm(5 << 20)
	upFile, _, err := r.FormFile("file-avatar")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer upFile.Close()

	targetFile, err := os.OpenFile("./assets/images/img-avatar.png", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, upFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)

	api.dashboardView(w, r)
}
