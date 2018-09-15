package main

import (
	"html/template"
    "fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

const dir = "tmpl/"

var templates = template.Must(template.ParseFiles(dir+"index.html", dir+"edit.html", dir+"view.html"))

func (p *Page) save() error {
	filename := "content/" + p.Title
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func load(title string) (*Page, error) {
	filename := "content/" + title
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return &Page{Title: "Error loading file"}, nil
	}
	return &Page{Title: title, Body: body}, nil
}

//common utility

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathValidator := regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
		m := pathValidator.FindStringSubmatch(r.URL.Path)
		if m == nil {
			// requst string is incorrect. What to do?
            w.Write([]byte(fmt.Sprintf("m: %v", m)))
        return 
		}
		fn(w, r, m[2])
	}
}

func renderTemplate(w http.ResponseWriter, title string, p *Page) {
	err := templates.ExecuteTemplate(w, title, p)
	if err != nil {
		// this is last call page can not be displayed do log and redirect back
		// todo: find how to get address of page where request was made from
	}

}
func handler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index", nil)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := load(title)
	if err != nil {
		// todo: write function to redirect to error page
		http.NotFound(w, r)
	}
	renderTemplate(w, "view", p)
}


// func editHandler(w http.ResponseWriter, r *http.Request, title, string){}
// func saveHandler(w http.ResponseWriter, r *http.Request, title, string){}

func main() {
	log.Println("Start Server")
    mux := http.NewServeMux()
	vh := makeHandler(viewHandler)
	mux.Handle("/",http.HandlerFunc( handler))
	//http.HandleFunc("/view/", viewH)
	mux.Handle("/view/", vh)
	// http.HandleFunc("/edit/", makeHandler(editHandler))
	// http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8989", mux)
}
