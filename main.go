package main

import (
	"io/ioutil"
	"os"

	"net/http"
	"regexp"
	"text/template"
)

var templates = template.Must(template.ParseFiles("layout.html"))
var validPath = regexp.MustCompile("^/(resume|portfolio|photos)|^/$|^$")

func makeLayout(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	} else {
		if m[0] == "" || m[0] == "/" {
			renderTemplate(w, "views/homepage.html")
		} else {
			renderTemplate(w, "views/"+m[1]+".html")
		}
	}
}

func renderTemplate(w http.ResponseWriter, filename string) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type content struct {
		Body string
	}
	viewData := content{
		Body: string(body),
	}
	err = templates.ExecuteTemplate(w, "layout.html", viewData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", makeLayout)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
