package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// File struct holds the basic structure of a file being returned
// it maybe a dir or a normal file
type File struct {
	IsDir   bool      `json:"isDir"`
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
	Mode    string    `json:"mode"`
}

func main() {
	portPtr := flag.Int("port", 7373, "port number")
	flag.Parse()

	portStr := fmt.Sprintf(":%d", *portPtr)

	router := mux.NewRouter()

	router.HandleFunc("/get", get).Methods("POST")

	log.Println(fmt.Sprintf("Starting server on port%s", portStr))
	listener, err := net.Listen("tcp", portStr)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	go http.Serve(listener, router)

	log.Println("Server is up now.")

	<-done
}

func setResponseTypeToJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func readDir(path string) []File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	files, _ := file.Readdir(0)

	var list []File

	for _, f := range files {
		file := File{IsDir: f.IsDir(), Name: f.Name(), Size: f.Size(), ModTime: f.ModTime(), Mode: f.Mode().String()}

		list = append(list, file)
	}

	return list
}

func get(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		Path string `json:"path"`
	}

	var reqBody ReqBody

	json.NewDecoder(r.Body).Decode(&reqBody)

	log.Println(reqBody.Path)

	dirs := readDir(reqBody.Path)
	log.Println(dirs)

	setResponseTypeToJSON(w)
	json.NewEncoder(w).Encode(dirs)
}
