package main

import (
  "io/ioutil"
  "net/http"
  "regexp"
  "os"
  "fmt"
  "strings"
  "log"
)

const (
	DEFAULT_PORT = ":8080"
)

var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")
var jsFileValidator = regexp.MustCompile(".+js$")
var cssFileValidator = regexp.MustCompile(".+css$")

//Page Request  data structure
type PageRequest struct {
  Method string
  Path string
  Filename string
}

func (pr *PageRequest)  New(path string, filename string, method string) {
  pr.Path = path
  pr.Filename = filename
  pr.Method = method

}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
		log.Println(path)
    if path == "/" {
      path = "/index.html"
    }
    filename := "public"+path
    splitPath :=  strings.Split(path,"/")
    if jsFileValidator.MatchString(splitPath[len(splitPath)-1]) {
      w.Header().Set("Content-Type", "text/javascript")
    } else if cssFileValidator.MatchString(splitPath[len(splitPath)-1]) {
      w.Header().Set("Content-Type", "text/css")
    } else {
		}
    body, err := ioutil.ReadFile(filename)
    if err != nil {
      
      errorHandler(w, r, err)
      return
    }
    fmt.Fprintf(w, string(body))
}

func errorHandler(w http.ResponseWriter, r *http.Request, err error) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  switch err.(type) {
  case *os.PathError:
    body,_ := ioutil.ReadFile("error/404.html")
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "%s",string(body))
  default:
    http.Error(w,"Fatal Error", http.StatusInternalServerError)
  }
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    fn(w, r)
  }
}

func setupServer(port string) {
	s := &http.Server{
		Addr:           port,
		Handler:        makeHandler(defaultHandler),
	}
	log.Printf("Starting server at port %s", port[1:])
  err := s.ListenAndServe()
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}

func main() {

	port := DEFAULT_PORT
	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}
	setupServer(port)
}
