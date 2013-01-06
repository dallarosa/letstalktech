package main

import (
  "io/ioutil"
  "net/http"
  "regexp"
//  "errors"
  "fmt"
  "strings"
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
    if path == "/" {
      path = "/index.html"
    }
    filename := "public"+path
    splitPath :=  strings.Split(path,"/")
    fmt.Println(splitPath[len(splitPath)-1])
    if jsFileValidator.MatchString(splitPath[len(splitPath)-1]) {
      w.Header().Set("Content-Type", "text/javascript")
      fmt.Println("It's a JS file")
    } else if cssFileValidator.MatchString(splitPath[len(splitPath)-1]) {
      w.Header().Set("Content-Type", "text/css")
      fmt.Println("It's a CSS file")
    }
    body, err := ioutil.ReadFile(filename)
    if err != nil {
      http.Error(w,err.Error(), http.StatusNotFound)
      return
    }
    fmt.Fprintf(w, string(body))
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Make Handler: %s\n", r.URL.Path)
    fn(w, r)
  }
}

func main() {
  http.HandleFunc("/", makeHandler(defaultHandler))
  http.ListenAndServe(":8080", nil)
}
