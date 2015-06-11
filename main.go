package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

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
	http.ListenAndServe(":8080", nil)
}
