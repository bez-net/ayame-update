package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Set of media files for service
type MediaSet struct {
	SrcDir   string   `json:"src_dir,omitempty"`
	DstDir   string   `json:"dst_dir,omitempty"`
	SrcBase  string   `json:"src_base,omitempty"`
	DstBase  string   `json:"dst_base,omitempty"`
	SrcName  string   `json:"src_name,omitempty"`
	DstName  string   `json:"dst_name,omitempty"`
	DstDesc  string   `json:"dst_desc,omitempty"`
	DstFiles []string `json:"dst_files,omitempty"`
}

// Stringer for MediaSet
func (m *MediaSet) String() string {
	return fmt.Sprintf("MediaSet> SrcDir=%s, DstDir=%s, SrcName=%s, DstName=%s, DstDesc=%s",
		m.SrcDir, m.DstDir, m.SrcName, m.DstName, m.DstDesc)
}

// Convert JSON output for the struct
func (m *MediaSet) JSONString() string {
	jbyte, _ := json.Marshal(m)
	return string(jbyte)
}

// Convert JSON Indent output for the struct
func (m *MediaSet) JSONIndentString() string {
	jbyte, _ := json.MarshalIndent(m, "", "  ")
	return string(jbyte)
}

// Send struct as the response in JSON
func (m *MediaSet) SendJSON(w http.ResponseWriter) {
	jbyte, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jbyte)
}

// --------------------------------------------------------------------------------------
// Handler for Uploading and Transcoding
func uploadHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path, r.RemoteAddr, r.Method, r.Header.Get("Content-Type"))
	defer log.Printf("uploadHandler exit")

	// parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files.
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("ParseMultipartForm error:", err)
		return
	}

	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		log.Println("FormFile error:", err)
		return
	}
	defer file.Close()

	// log.Println(handler)

	basename := time.Now().Format("D20060102T150405")

	// create a temp file within our upload directory that follows a particular naming pattern
	tempFile, err := ioutil.TempFile("asset/upload", "COBOT-"+basename+"-R*"+filepath.Ext(handler.Filename))
	if err != nil {
		log.Println("TempFile error:", err)
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
	fmt.Fprintf(w, "Successfully uploaded the file and being processed it.\n")
	log.Printf("%s => %s", handler.Filename, tempFile.Name())

	// prepare a media set for the upload file
	mset := &MediaSet{}
	mset.SrcDir = "asset/upload/"
	mset.SrcBase = ""
	mset.SrcName = filepath.Base(tempFile.Name())
	mset.DstDir = "asset/record/"
	mset.DstBase = basename + "/" // time.Now().Format("D20060102T150405/")
	mset.DstName = "COBOT-" + basename + "-U" + getUUIDString()
	mset.DstDesc = "ffmpeg (libx264/aac) mp4/mpv/jpg/vtt"
	log.Println("MediaSet>\n", mset.JSONIndentString())

	// do post media processing in background
	go postMediaProcessing(mset)
}

// Postprocessing for the video file uploaded
func postMediaProcessing(mset *MediaSet) (err error) {
	defer log.Printf("postMediaProcessing Done")

	err = getMediaInfo(mset)
	if err != nil {
		log.Println("getMediaInfo error:", err)
		return
	}
	// log.Println("getMediaInfo:", "Done")

	err = makeMediaSet(mset)
	if err != nil {
		log.Println("makeMediaSet error:", err)
		return
	}
	// log.Println("makeMediaSet:", "Done")

	// remove the upload file after sucessful record
	err = os.Remove(mset.SrcDir + mset.SrcName)
	if err != nil {
		log.Println("Remove error:", err)
		return
	}
	return
}

// Make a set of media files for a video
func getMediaInfo(mset *MediaSet) (err error) {
	// check mediainfo command if it is executable
	_, err = exec.LookPath("mediainfo")
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("mediainfo:", path)

	// Get meta information for the media file
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("mediainfo", mset.SrcDir+mset.SrcName)
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
func makeMediaSet(mset *MediaSet) (err error) {
	_, err = exec.LookPath("ffmpeg")
	if err != nil {
		log.Printf("LookPath error: %s", err)
		return
	}
	// log.Println("ffmpeg:", path)

	// generate related files for the input video
	os.MkdirAll(mset.DstDir+mset.DstBase, os.ModePerm)

	inPart := mset.SrcDir + mset.SrcName
	outPart := mset.DstDir + mset.DstBase + mset.DstName
	log.Println(inPart, "=>", outPart)

	cmdStr := fmt.Sprintf("ffmpeg -loglevel error -stats -y")
	cmdStr += fmt.Sprintf(" -i %s", inPart)

	// if you want to use libfdk_aac, check its support in ffmpeg -codecs / -encoders
	mp4Opt := `-vcodec libx264 -vf "scale=1280:720" -acodec aac -movflags faststart -f mp4`
	mp4Part := changePathExtention(outPart, ".mp4")
	cmdStr += fmt.Sprintf(" %s %s", mp4Opt, mp4Part)

	// consider use middle(480:270) if the size(320:180) is small
	mpvOpt := `-vcodec libx264 -r 10 -vf "scale=480:270" -an -movflags faststart -f mp4`
	mpvPart := changePathExtention(outPart, ".mpv")
	cmdStr += fmt.Sprintf(" %s %s", mpvOpt, mpvPart)

	jpgOpt := `-ss 00:00:01.000 -frames:v 1 -vf "scale=640:360"`
	jpgPart := changePathExtention(outPart, ".jpg")
	cmdStr += fmt.Sprintf(" %s %s", jpgOpt, jpgPart)

	// gifOpt := `-r 10 -vf "scale=320:180" -loop 0`
	// gifPart := changePathExtention(outPart, ".gif")
	// cmdStr += fmt.Sprintf("%s %s", gifOpt, gifPart)

	// webpOpt := `-r 10 -vf "scale=320:180" -loop 0`
	// webpPart := changePathExtention(outPart, ".webp")
	// cmdStr += fmt.Sprintf(" %s %s", webpOpt, webpPart)

	// webmOpt := `-r 10 -vf "scale=320:180" -an`
	// webmPart := changePathExtention(outPart, ".webm")
	// cmdStr += fmt.Sprintf(" %s %s", webmOpt, webmPart)

	log.Println(cmdStr)
	cmd := exec.Command("sh", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("CombinedOutput error: %s", err)
		return
	}
	log.Println(string(out))

	vttPart := changePathExtention(outPart, ".vtt")
	err = makeSubtitleFile(vttPart)
	if err != nil {
		log.Printf("makeSubtitleFile err: %s", err)
		return
	}

	err = recordMediaSetFiles(mset)
	if err != nil {
		log.Printf("recordMediaSetFiles err: %s", err)
		return
	}
	log.Println("MediaSet>\n", mset.JSONIndentString())
	return
}

// Change extension of the filename with another one
func changePathExtention(fpath, fext string) (str string) {
	ext := filepath.Ext(fpath)
	str = fpath[0:len(fpath)-len(ext)] + fext
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

const sampleTitle = `WEBVTT - Sample text for testing
NOTE
This is comment part. It is not displayed on screen.

00:00:01.000 --> 00:00:05.000 line:-3
<b>This is sample subtitle text.</b>
HTML style tag can be used.
`

func makeSubtitleFile(fname string) (err error) {
	f, err := os.Create(fname)
	if err != nil {
		log.Printf("Create error: %s", err)
		return
	}
	defer f.Close()

	n, err := f.WriteString(sampleTitle)
	if n < len(sampleTitle) || err != nil {
		log.Printf("WriteString error: %s, %d", err, n)
		return
	}
	f.Sync()
	return
}

func recordMediaSetFiles(mset *MediaSet) (err error) {
	d, err := os.Open(mset.DstDir + mset.DstBase)
	if err != nil {
		log.Printf("Open error: %s", err)
		return
	}
	defer d.Close()
	files, err := d.Readdir(-1)
	if err != nil {
		log.Printf("Readdir error: %s", err)
		return
	}
	for _, f := range files {
		if f.Mode().IsRegular() && f.Size() > 0 {
			mset.DstFiles = append(mset.DstFiles, f.Name())
		}
	}
	return
}
