package thereg

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"
)

var templates = template.Must(template.ParseFiles("tmpl/index.html", "tmpl/edit.html", "tmpl/view.html"))
var validPath = regexp.MustCompile("^/(account|node)/([a-zA-Z0-9]+)$")
var validData = regexp.MustCompile("^data/([a-zA-Z0-9]+)\\.txt$")

// func getID(w http.ResponseWriter, r *http.Request) (string, error) {
// 	m := validPath.FindStringSubmatch(r.URL.Path)
// 	if m == nil {
// 		http.NotFound(w, r)
// 		return "", errors.New("Invalid Node ID")
// 	}
// 	return m[2], nil // The ID is the second subexpression.
// }

// func renderTemplate(w http.ResponseWriter, tmpl string, n *Node) {
// 	err := templates.ExecuteTemplate(w, tmpl+".html", n)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// func viewHandler(w http.ResponseWriter, r *http.Request, id string) {
// 	fmt.Println(time.Now().UTC().UnixNano(), "Viewing", id)
// 	n, err := loadNode(id)
// 	if err != nil {
// 		http.Redirect(w, r, "/edit/"+id, http.StatusFound)
// 		return
// 	}

// 	renderTemplate(w, "view", n)
// }

// func editHandler(w http.ResponseWriter, r *http.Request, id string) {
// 	fmt.Println(time.Now().UTC().UnixNano(), "Editing", id)
// 	n, err := loadNode(id)
// 	if err != nil {
// 		n = &Node{ID: id}
// 	}

// 	renderTemplate(w, "edit", n)
// }

// func saveHandler(w http.ResponseWriter, r *http.Request, id string) {
// 	fmt.Println(time.Now().UTC().UnixNano(), "Saving", id)
// 	URL := r.FormValue("URL")
// 	n := &Node{ID: id, URL: string(URL)}
// 	err := n.save()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	http.Redirect(w, r, "/view/"+id, http.StatusFound)
// }

// func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Here we will extract the page id from the Request,
// 		// and call the provided handler 'fn'
// 		m := validPath.FindStringSubmatch(r.URL.Path)
// 		if m == nil {
// 			http.NotFound(w, r)
// 			return
// 		}
// 		fn(w, r, m[2])
// 	}
// }

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC().UnixNano(), "Index")

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

func nodesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC().UnixNano(), "/node")

	w.Header().Set("Content-Type", "application/json")

	hostname := r.FormValue("hostname")
	port := r.FormValue("port")
	params := Node{Hostname: hostname, Port: port}

	var err error

	nodes := NodeList{}
	switch r.Method {
	case http.MethodGet:
		// GET /node
		res := nodes.Index()
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodPost:
		// POST /node
		res := nodes.Create(params)
		js, _ := json.Marshal(res)
		w.Write(js)
	default:
		// Give an error message.
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func nodeHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	id := m[2]

	fmt.Println(time.Now().UTC().UnixNano(), "/node/:id", id)

	name := r.FormValue("name")
	hostname := r.FormValue("hostname")
	port := r.FormValue("port")
	status := r.FormValue("status")
	params := Node{Name: name, Hostname: hostname, Port: port, Status: status}

	var err error

	nodes := NodeList{}
	switch r.Method {
	case http.MethodGet:
		// GET /node/:id
		res := nodes.Read(id)
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodPatch:
		// PATCH /node/:id
		res := nodes.Update(id, params)
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodDelete:
		// DELETE /node/:id
		res := nodes.Delete(id)
		js, _ := json.Marshal(res)
		w.Write(js)
	default:
		// Give an error message.
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// POST /account
// GET /account
// PATCH /account
// DELETE /account
func accountHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC().UnixNano(), "/account")

	w.Header().Set("Content-Type", "application/json")

	email := r.FormValue("email")
	username := r.FormValue("username")
	params := Account{Email: email, Username: username}

	var err error

	accounts := AccountList{}
	switch r.Method {
	case http.MethodPost:
		// POST /account
		res := accounts.Create()
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodGet:
		// GET /account
		res := accounts.Read()
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodPatch:
		// PATCH /account
		res := accounts.Update(params)
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodDelete:
		// DELETE /account
		res := accounts.Delete()
		js, _ := json.Marshal(res)
		w.Write(js)
	default:
		// Give an error message.
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func accountConfirmHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC().UnixNano(), "/account")

	accounts := AccountList{}
	account := accounts.Read()

	js, err := json.Marshal(account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Serve runs the http server
func Serve() {

	http.HandleFunc("/", indexHandler)

	// GET /node
	// POST /node
	http.HandleFunc("/node", nodesHandler)
	// GET /node/:id
	// PATCH /node/:id
	// DELETE /node/:id
	http.HandleFunc("/node/", nodeHandler)

	// POST /account
	// GET /account
	// PATCH /account
	// DELETE /account
	http.HandleFunc("/account", accountHandler)
	// GET /account/confirm/:token
	http.HandleFunc("/account/confirm", accountConfirmHandler)

	fmt.Println("Server running at", "localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
