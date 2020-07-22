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
	uploadVideos(w, r)
}

func uploadVideo(w http.ResponseWriter, r *http.Request) {
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
		log.Fatal(err)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
		return
	}

	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "File uploaded\n")
}

func uploadVideos(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	files := r.MultipartForm.File["myFiles"]

	for _, file := range files {
		f, _ := file.Open()

		fmt.Printf("Uploaded File: %+v\n", file.Filename)
		fmt.Printf("File Size: %+v\n", file.Size)
		fmt.Printf("MIME Header: %+v\n", file.Header)
		fmt.Println("----------------------------------")

		tempFile, err := ioutil.TempFile("videos", "upload-*.mp4")
		if err != nil {
			log.Fatal(err)
			return
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
			return
		}

		tempFile.Write(fileBytes)
		fmt.Fprintf(w, "Uploaded file: %+v\n", file.Filename)
	}

}
