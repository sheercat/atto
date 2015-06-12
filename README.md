# atto

super handy http server. with upload form, basic auth.

## install

    go get github.com/sheercat/atto

## Getting Started
```
   $ atto

   $ atto --port=9090

   $ atto --user=basic --pass=auth
```
   then, atto serve current dir files to default port 8080.

   If you designate params '?upload' to url then atto serve html with file upload form.
   and you can upload file to current dir.

## License

atto licensed under the MIT



