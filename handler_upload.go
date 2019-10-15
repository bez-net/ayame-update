package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func uploadHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Printf(r.URL.Path)
	// send a response as a upload page
	page, err := ioutil.ReadFile("./static/upload.html")
	if err != nil {
		log.Printf("ReadFile error: %s", err)
		return
	}
	fmt.Fprintf(w, string(page))
	// Parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Printf("FormFile error: %s", err)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temp file within our upload directory that follows a particular naming pattern
	tempFile, err := ioutil.TempFile("upload", "upload-*.png")
	if err != nil {
		log.Printf("TempFile error: %s", err)
		return
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("ReadAll error: %s", err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully uploaded file\n")
}
