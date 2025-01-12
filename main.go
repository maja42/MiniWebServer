package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
)

const disableCaching = true

// Respond to the URL /home with an html home page
func FileHandler(response http.ResponseWriter, request *http.Request) {
	var requestPath = request.URL.Path
	log.Printf("Requested %q", requestPath)

	if requestPath[0] != byte('/') {
		panic("Unable to serve file: unexpected url")
	}

	var filePath = "." + requestPath

	resHeader := response.Header()
	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	if contentType != "" {
		resHeader.Set("Content-Type", contentType)
	}

	webpage, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Responding with 404 (file not found)")
		http.Error(response, fmt.Sprintf("Failed to serve file: %v", err), 404)
		return
	}

	if disableCaching {
		resHeader.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		resHeader.Set("Pragma", "no-cache")
		resHeader.Set("Expires", "0")
	}

	log.Printf("Serving %q (%d bytes)", filePath, len(webpage))
	fmt.Fprint(response, string(webpage))
}

func main() {
	port := 8088
	portstring := strconv.Itoa(port)

	// Register request handlers for two URL patterns.
	// (The docs are unclear on what a 'pattern' is,
	// but seems be the start of the URL, ending in a /).
	// See gorilla/mux for a more powerful matching system.
	// Note that the "/" pattern matches all request URLs.
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(FileHandler))

	log.Printf("This webserver serves all files in the current working directory.")
	log.Printf("Listening on port " + portstring + " ... ")
	err := http.ListenAndServe(":"+portstring, mux)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
