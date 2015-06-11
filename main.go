package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

var PortNumber *string = flag.String("port", "8080", "port number.")

func rootHandler(w http.ResponseWriter, r *http.Request) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	http.ServeFile(w, r, path.Join(dir, r.URL.Path))
}

func main() {
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":"+*PortNumber, nil))
}
