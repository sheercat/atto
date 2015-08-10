package main

import (
	"flag"
	"fmt"
	"github.com/k0kubun/pp"
	//_ "github.com/motemen/go-loghttp/global"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

var portNumber = flag.String("port", "8080", "port number.")
var basicAuthUser = flag.String("user", "", "basic auth user name")
var basicAuthPass = flag.String("pass", "", "basic auth user pass")
var debugFlag = flag.Bool("debug", false, "debug")

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

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	// log.Println(username, password, ok)
	if ok == false {
		return false
	}
	return username == *basicAuthUser && password == *basicAuthPass
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if *debugFlag {
		pp.Println(r)
	}
	if *basicAuthUser != "" && *basicAuthPass != "" {
		log.Println(*basicAuthUser)
		if checkAuth(w, r) == false {
			w.Header().Set("WWW-Authenticate", `Basic realm="Atto"`)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
			return
		}
	}

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
