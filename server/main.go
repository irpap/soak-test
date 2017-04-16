package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
)

var uploadDir string

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	os.MkdirAll(uploadDir, 0600)
	vars := mux.Vars(r)
	filename := vars["key"]

	out, err := os.Create(path.Join(uploadDir, filename))
	if err != nil {
		log.Fatal(err)
	}
	if _, err := io.Copy(out, r.Body); err != nil {
		log.Fatal(err)
	}
	r.Body.Close()
	return
}

func main() {
	flag.StringVar(&uploadDir, "uploadDir", "/tmp/uploads/", "Directory to store uploaded files in")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/{key:[[:alnum:]\\._-]+}", uploadHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))

}
