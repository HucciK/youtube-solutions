package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func (h Handler) ClientUpdate(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/clientUpdate/"):]

	if path == "RELEASES" {
		releases, err := os.ReadFile("../resources/RELEASES")
		if err != nil {
			fmt.Println(err)
		}

		w.Write(releases)
		return
	}

	update, err := os.ReadFile(fmt.Sprintf("../resources/%s", path))
	if err != nil {
		fmt.Println(err)
	}

	w.Write(update)
}

func (h Handler) DevUpdate(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/devUpdate/"):]

	if path == "RELEASES" {
		releases, err := os.ReadFile("../dev_resources/RELEASES")
		if err != nil {
			fmt.Println(err)
		}

		w.Write(releases)
		return
	}

	update, err := os.ReadFile(fmt.Sprintf("../dev_resources/%s", path))
	if err != nil {
		fmt.Println(err)
	}

	if _, err = w.Write(update); err != nil {
		w.WriteHeader(500)
	}
}

func (h Handler) GetChangeLog(w http.ResponseWriter, r *http.Request) {
	theme := r.URL.Query().Get("theme")

	changeLog, err := template.ParseFiles(fmt.Sprintf("../resources/html/changelog_%s.html", theme))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if err = changeLog.Execute(w, nil); err != nil {
		w.WriteHeader(500)
	}
}

func (h Handler) GetCSS(w http.ResponseWriter, r *http.Request) {

	filename := r.URL.Path[len("/css/"):]

	file, err := os.ReadFile(fmt.Sprintf("../resources/css/%s", filename))
	if err != nil {
		return
	}
	w.Header().Add("Content-Type", "text/css")

	if _, err = w.Write(file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h Handler) GetJS(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path[len("/scripts/"):]

	file, err := os.ReadFile(fmt.Sprintf("../resources/scripts/%s", filename))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
