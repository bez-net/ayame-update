package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func fetchHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path, r.RemoteAddr, r.Method, r.Header.Get("Content-Type"))
	defer log.Printf("fetchHandler exit")

	// Parse our multipart form, 10 << 20 specifies a maximum upload of 32 MB files.
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Printf("ParseMultipartForm error: %s", err)
		return
	}

	// for debugging, check the elements of multipart form
	for k := range r.MultipartForm.File {
		log.Println(k)
	}

	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		log.Printf("FormFile error: %s", err)
		return
	}
	defer file.Close()

	// Create a temp file within our upload directory that follows a particular naming pattern
	tempFile, err := ioutil.TempFile("upload", "cojam-*"+filepath.Ext(handler.Filename))
	if err != nil {
		log.Printf("TempFile error: %s", err)
		return
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("ReadAll error: %s", err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully uploaded file\n")
	log.Printf("%s is stored to %s", handler.Filename, tempFile.Name())
}
