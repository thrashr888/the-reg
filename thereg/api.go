package thereg

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("tmpl/index.html"))
var validPath = regexp.MustCompile("^/api/(account|node)/([a-zA-Z0-9\\-_]+)$")
var validConfirmPath = regexp.MustCompile("^/api/account/confirm/([a-zA-Z0-9\\-_]+)$")

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", map[string]string{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func nodesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	account := DBGetAccountByToken(string(r.Header["Auth-Token"][0]))

	decoder := json.NewDecoder(r.Body)
	var params Node
	err := decoder.Decode(&params)

	nodes := NodeList{}
	switch r.Method {
	case http.MethodGet:
		// GET /api/node
		res := nodes.Index(account)
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodPost:
		// POST /api/node
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
	w.Header().Set("Content-Type", "application/json")

	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	id := m[2]

	account := DBGetAccountByToken(string(r.Header["Auth-Token"][0]))
	if account.ID == "" {
		http.Error(w, "Account not found", http.StatusForbidden)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var params Node
	err := decoder.Decode(&params)

	nodes := NodeList{}
	switch r.Method {
	case http.MethodGet:
		// GET /api/node/:id
		res := nodes.Read(id)
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodPatch:
		// PATCH /api/node/:id
		res := nodes.Update(account, id, params)
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodDelete:
		// DELETE /api/node/:id
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

// POST /api/account
// GET /api/account
// PATCH /api/account
// DELETE /api/account
func accountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	account := DBGetAccountByToken(string(r.Header["Auth-Token"][0]))

	decoder := json.NewDecoder(r.Body)
	var params Account
	err := decoder.Decode(&params)

	accounts := AccountList{}
	switch r.Method {
	case http.MethodPost:
		// POST /api/account
		res := accounts.Create()
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodGet:
		// GET /api/account
		res := accounts.Read(account)
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodPatch:
		// PATCH /api/account
		res := accounts.Update(account, params)
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodDelete:
		// DELETE /api/account
		res := accounts.Delete(account)
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
	w.Header().Set("Content-Type", "application/json")

	m := validConfirmPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	confirmToken := m[2]

	accounts := AccountList{}
	res := accounts.Confirm(confirmToken)
	js, _ := json.Marshal(res)
	w.Write(js)
}

func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("remote-addr=%s; method=%s; host=%s; uri=%s; content-type=%s;", r.RemoteAddr, r.Method, r.Host, r.RequestURI, r.Header["Content-Type"])
		next.ServeHTTP(w, r)
	}
}

type middleware func(next http.HandlerFunc) http.HandlerFunc

func chainMiddleware(mw ...middleware) middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}

// ServeAPI runs the http/API server
func ServeAPI(port string) {
	lt := chainMiddleware(withLogging)

	// GET /
	http.HandleFunc("/", lt(indexHandler))

	// GET /api/node
	// POST /api/node
	http.HandleFunc("/api/node", lt(nodesHandler))
	// GET /api/node/:id
	// PATCH /api/node/:id
	// DELETE /api/node/:id
	http.HandleFunc("/api/node/", lt(nodeHandler))

	// POST /api/account
	// GET /api/account
	// PATCH /api/account
	// DELETE /api/account
	http.HandleFunc("/api/account", lt(accountHandler))
	// GET /api/account/confirm/:token
	http.HandleFunc("/api/account/confirm/", lt(accountConfirmHandler))

	host := fmt.Sprintf(":%s", port)
	log.Println("API Server running at", host)
	log.Fatal(http.ListenAndServe(host, nil))
}
