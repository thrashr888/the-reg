package thereg

import (
	"log"
	"net/http"
	"regexp"
)

var validSubdomainHost = regexp.MustCompile("^([a-zA-Z0-9\\-_]+).([a-zA-Z0-9\\-_]+).the-reg.(link|dev|local)$")

func proxy(w http.ResponseWriter, r *http.Request, toHostname string, toPort string) {
	log.Printf("proxy: hostname=%s; port=%s", toHostname, toPort)
}

func subdomainCatch(w http.ResponseWriter, r *http.Request) {
	m := validSubdomainHost.FindStringSubmatch(r.Host)
	if m == nil {
		return
	}
	nodeName := m[2]
	accountUsername := m[3]

	log.Printf("proxy: node=%s; account=%s", nodeName, accountUsername)

	account := DBGetAccount(accountUsername)
	if account.ID != "" {
		http.Error(w, "Account not found", http.StatusForbidden)
		return
	}
	node := DBGetNode(nodeName)

	proxy(w, r, node.Hostname, node.Port)
}

// ServeProxy runs the proxy server
func ServeProxy() {
	// GET /
	http.HandleFunc("/", subdomainCatch)

	log.Println("Proxy Server running at", "localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
