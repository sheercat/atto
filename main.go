package main

import (
	"flag"
	"fmt"
	// _ "github.com/k0kubun/pp"
	//_ "github.com/motemen/go-loghttp/global"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

var portNumber = flag.String("port", "8080", "port number.")

func uploadForm(w http.ResponseWriter) {
	fmt.Fprintln(w, `
<html><head><title>Go upload</title></head>
<body>
<form action="./" method="post" enctype="multipart/form-data">
<label for="file">Filename:</label>
<input type="file" name="file" id="file">
<input type="submit" name="submit" value="Submit">
</form>
</body>
</html>
`)
}

func uploadFile(w http.ResponseWriter, r *http.Request) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()

	targetFile := path.Join(dir, r.URL.Path, header.Filename)
	out, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}
	return nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// pp.Println(r)
	log.Println(r.Method, r.RequestURI)
	if r.URL.RawQuery == "upload" {
		uploadForm(w)
		return
	}
	if r.Method == "POST" {
		err := uploadFile(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
		} else {
			http.Redirect(w, r, r.URL.Path, 301)
		}
		return
	}
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	http.ServeFile(w, r, path.Join(dir, r.URL.Path))
}

func main() {
	flag.Parse()
	http.HandleFunc("/", rootHandler)
	log.Println("listen:" + *portNumber)
	log.Fatal(http.ListenAndServe(":"+*portNumber, nil))
}
