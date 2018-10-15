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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC().UnixNano(), "Index")

	err := templates.ExecuteTemplate(w, "index.html", map[string]string{"title": "Home"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func nodesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().UTC().UnixNano(), "/node")

	w.Header().Set("Content-Type", "application/json")

	account := DBGetAccountByToken(string(r.Header["Auth-Token"][0]))

	decoder := json.NewDecoder(r.Body)
	var params Node
	err := decoder.Decode(&params)

	nodes := NodeList{}
	switch r.Method {
	case http.MethodGet:
		// GET /node
		res := nodes.Index(account)
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodPost:
		// POST /node
		res := nodes.Create(account, params)
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
	fmt.Println(time.Now().UTC().UnixNano(), "/node/:id")

	w.Header().Set("Content-Type", "application/json")

	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	id := m[2]

	fmt.Println(time.Now().UTC().UnixNano(), "/node/", id)

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
