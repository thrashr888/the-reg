package thereg

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("tmpl/index.html", "tmpl/edit.html", "tmpl/view.html"))
var validPath = regexp.MustCompile("^/(account|node)/([a-zA-Z0-9\\-_]+)$")
var validConfirmPath = regexp.MustCompile("^/account/confirm/([a-zA-Z0-9\\-_]+)$")
var validSubdomainHost = regexp.MustCompile("^([a-zA-Z0-9\\-_]+).([a-zA-Z0-9\\-_]+).the-reg.(link|dev|local)$")

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", map[string]string{"title": "Home"})
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
		// GET /node/:id
		res := nodes.Read(id)
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodPatch:
		// PATCH /node/:id
		res := nodes.Update(account, id, params)
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
	w.Header().Set("Content-Type", "application/json")

	account := DBGetAccountByToken(string(r.Header["Auth-Token"][0]))

	decoder := json.NewDecoder(r.Body)
	var params Account
	err := decoder.Decode(&params)

	accounts := AccountList{}
	switch r.Method {
	case http.MethodPost:
		// POST /account
		res := accounts.Create()
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodGet:
		// GET /account
		res := accounts.Read(account)
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodPatch:
		// PATCH /account
		res := accounts.Update(account, params)
		js, _ := json.Marshal(res)
		w.Write(js)
	case http.MethodDelete:
		// DELETE /account
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

func proxyHandler(node Node, w http.ResponseWriter, r *http.Request) {
	// TODO add the proxy
}

func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("remote-addr=%s; method=%s; host=%s; uri=%s; content-type=%s;", r.RemoteAddr, r.Method, r.Host, r.RequestURI, r.Header["Content-Type"])
		next.ServeHTTP(w, r)
	}
}

func withSubdomainCatch(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		m := validSubdomainHost.FindStringSubmatch(r.Host)
		if m == nil {
			next.ServeHTTP(w, r)
			return
		}
		nodeName := m[2]
		accountUsername := m[3]

		log.Printf("node=%s; account=%s", nodeName, accountUsername)

		account := DBGetAccount(accountUsername)
		if account.ID != "" {
			http.Error(w, "Account not found", http.StatusForbidden)
			return
		}
		node := DBGetNode(nodeName)

		proxyHandler(node, w, r)

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

// Serve runs the http server
func Serve() {
	lt := chainMiddleware(withLogging, withSubdomainCatch)

	// GET /
	http.HandleFunc("/", lt(indexHandler))

	// GET /node
	// POST /node
	http.HandleFunc("/node", lt(nodesHandler))
	// GET /node/:id
	// PATCH /node/:id
	// DELETE /node/:id
	http.HandleFunc("/node/", lt(nodeHandler))

	// POST /account
	// GET /account
	// PATCH /account
	// DELETE /account
	http.HandleFunc("/account", lt(accountHandler))
	// GET /account/confirm/:token
	http.HandleFunc("/account/confirm/", lt(accountConfirmHandler))

	log.Println("Server running at", "localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
