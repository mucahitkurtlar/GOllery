package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob("templates/*.html"))

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	vs := http.FileServer(http.Dir("./videos/"))
	r.PathPrefix("/videos/").Handler(http.StripPrefix("/videos/", vs))

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/upload", uploadHandler)
	r.HandleFunc("/upload-video", uploadVideoHandler)
	r.HandleFunc("/watch", watchHandler)
	http.Handle("/", r)

	fmt.Println("Listening on port 8080")
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	templates.ExecuteTemplate(w, "index.html", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	templates.ExecuteTemplate(w, "upload.html", nil)
}

func watchHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	templates.ExecuteTemplate(w, "watch.html", nil)
}

func uploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	tempFile, err := ioutil.TempFile("videos", "upload-*.mp4")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
