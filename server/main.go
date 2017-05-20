package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
)

var uploadDir string

func createProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	directory := vars["dir"]
	os.RemoveAll(path.Join(uploadDir, directory))
	if err := os.Mkdir(path.Join(uploadDir, directory), os.FileMode(0755)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func uploadPicture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	directory := vars["dir"]
	filename := vars["picture"]
	if _, err := os.Stat(path.Join(uploadDir, directory)); os.IsNotExist(err) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	out, err := os.Create(path.Join(uploadDir, directory, filename))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, r.Body); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func main() {
	flag.StringVar(&uploadDir, "uploadDir", "/tmp/uploads/", "Directory to store uploaded files in")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/{dir:[[:alnum:]_-]+}/", createProfile).Methods("POST")
	r.HandleFunc("/{dir:[[:alnum:]_-]+}/{picture:[[:alnum:]_-]+}", uploadPicture).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(uploadDir)))

	log.Fatal(http.ListenAndServe(":8000", r))

}
