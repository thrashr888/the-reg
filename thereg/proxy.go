package thereg

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
)

// curl localhost:8081 -H"Host:redis.thrashr888.proxy.the-reg.local"
// curl localhost:8081 -H"Host:http.thrashr888.proxy.the-reg.local:8080"
// redis.thrashr888.proxy.the-reg.local
var validSubdomainHost = regexp.MustCompile("^([a-zA-Z0-9\\-_]+).([a-zA-Z0-9\\-_]+).proxy.the-reg.(link|dev|local)(:[0-9]{2,4})?$")

// ServeProxy runs the proxy server
func ServeProxy() {
	director := func(req *http.Request) {
		var requestHostname string
		if req.URL.Host != "" {
			requestHostname = req.URL.Host
		} else {
			requestHostname = req.Host
		}

		m := validSubdomainHost.FindStringSubmatch(requestHostname)
		if m == nil {
			log.Printf("no hostmatch: hostname=%s", requestHostname)
			return
		}
		nodeName := m[1]
		accountUsername := m[2]
		log.Printf("lookup: node=%s; account=%s", nodeName, accountUsername)

		account := DBGetAccount(accountUsername)
		if account.ID == "" {
			log.Println("Account not found", accountUsername, account.ID)
			return
		}
		node := DBGetNode(nodeName)

		originHostname := fmt.Sprintf("%s:%s", node.Hostname, node.Port)

		log.Printf("direct: req=%s; origin=%s", requestHostname, originHostname)

		req.Header.Add("X-Forwarded-Host", requestHostname)
		req.Header.Add("X-Origin-Host", originHostname)
		req.URL.Scheme = "http"
		req.URL.Host = originHostname
	}

	proxy := &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	port := "8081"
	log.Println("API Server running at", "localhost:", port)
	host := fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(host, nil))
}
