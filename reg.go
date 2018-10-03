package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// User represents... the user
type User struct {
	ID                string
	Email             string
	EmailConfirmToken string
	EmailConfirmed    bool
	Username          string
	IP                string
	Authtoken         string
}

// UserList is a collection of Users
type UserList struct {
	Users []User
}

// Node represents a single networked resource
type Node struct {
	ID       string
	Name     string
	URL      string
	Hostname string
	Port     string
	Status   string
	Public   bool
}

// NodeList is a collection of Nodes
type NodeList struct {
	Nodes []Node
}

func (p *Node) save() error {
	filename := "data/" + p.ID + ".txt"
	return ioutil.WriteFile(filename, []byte(p.URL), 0600)
}

var templates = template.Must(template.ParseFiles("tmpl/index.html", "tmpl/edit.html", "tmpl/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var validData = regexp.MustCompile("^data/([a-zA-Z0-9]+)\\.txt$")

func getID(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Node ID")
	}
	return m[2], nil // The ID is the second subexpression.
}

func loadNode(id string) (*Node, error) {
	filename := "data/" + id + ".txt"

	url, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Node{ID: id, URL: string(url)}, nil
}

func loadNodes() ([]Node, error) {
	files, err := filepath.Glob("./data/*.txt")
	if err != nil {
		log.Fatal(err)
	}

	nodes := []Node{}
	for _, f := range files {
		_, file := filepath.Split(f)
		id := strings.TrimSuffix(file, path.Ext(file))
		if err != nil {
			continue
		}
		nodes = append(nodes, Node{ID: id})
	}

	return nodes, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, n *Node) {
	err := templates.ExecuteTemplate(w, tmpl+".html", n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, id string) {
	fmt.Println("Viewing", id)
	n, err := loadNode(id)
	if err != nil {
		http.Redirect(w, r, "/edit/"+id, http.StatusFound)
		return
	}

	renderTemplate(w, "view", n)
}

func editHandler(w http.ResponseWriter, r *http.Request, id string) {
	fmt.Println("Editing", id)
	n, err := loadNode(id)
	if err != nil {
		n = &Node{ID: id}
	}

	renderTemplate(w, "edit", n)
}

func saveHandler(w http.ResponseWriter, r *http.Request, id string) {
	fmt.Println("Saving", id)
	URL := r.FormValue("URL")
	n := &Node{ID: id, URL: string(URL)}
	err := n.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+id, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Here we will extract the page id from the Request,
		// and call the provided handler 'fn'
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Index")

	n, err := loadNodes()
	if n == nil {
		http.NotFound(w, r)
		return
	}

	err = templates.ExecuteTemplate(w, "index.html", NodeList{Nodes: n})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	fmt.Println("Server running at", "localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
