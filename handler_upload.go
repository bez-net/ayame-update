package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func uploadHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Printf("%s, %s", r.URL.Path, r.RemoteAddr)
	defer log.Printf("uploadHandler exit")

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

	// fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	// fmt.Printf("File Size: %+v\n", handler.Size)
	// fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temp file within our upload directory that follows a particular naming pattern
	tempFile, err := ioutil.TempFile("upload", "cobot-*"+filepath.Ext(handler.Filename))
	if err != nil {
		log.Printf("TempFile error: %s", err)
		return
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("ReadAll error:", err)
		return
	}
	// write this byte array to our temporary file
	n, err := tempFile.Write(fileBytes)
	if int64(n) < handler.Size || err != nil {
		log.Println("Write error:", err, n)
		return
	}

	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully did upload the file and process it.\n")
	log.Printf("%s is stored to %s", handler.Filename, tempFile.Name())

	// Do post media processing in background
	go postMediaProcessing(tempFile)
}

func postMediaProcessing(mediaFile *os.File) (err error) {
	err = getMediaInfo(mediaFile)
	if err != nil {
		log.Println("getMediaInfo error:", err)
		return
	}
	log.Println("getMediaInfo")

	err = makeMediaSet(mediaFile)
	if err != nil {
		log.Println("makeMediaSet error:", err)
		return
	}
	log.Println("makeMediaSet")
	return
}

// Make a set of media files for a video
func getMediaInfo(mediaFile *os.File) (err error) {
	// check mediainfo command if it is executable
	path, err := exec.LookPath("mediainfo")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("mediainfo:", path)

	// Get meta information for the media file
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("mediainfo", mediaFile.Name())
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err, string(stderr.Bytes()))
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	log.Println(outStr, errStr)
	return
}

// Make a set of media files for a video
func makeMediaSet(mediaFile *os.File) (err error) {
	path, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ffmpeg:", path)

	// Get meta information for the media file
	inPart := mediaFile.Name()
	outPart := "record/sample.webp"
	log.Println(inPart, outPart)
	cmd := exec.Command("ffmpeg", "-y", "-i", inPart, outPart)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	log.Println(string(out))
	return
}

// Send a web page to the http client
func sendFilePage(w http.ResponseWriter, filename string) (err error) {
	page, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("ReadFile(%s) error: %s", filename, err)
		return
	}
	fmt.Fprintf(w, string(page))
	return
}
