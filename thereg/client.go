package thereg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var apiHost = "localhost:8080/api"

// GetNodes - GET /api/node
func GetNodes() NodeList {
	res, _ := apiGET("node", "")
	n := NodeList{}
	return n.ListFromJSONBody(res)
}

// CreateNode - POST /api/node
func CreateNode(node Node) Node {
	params, _ := json.Marshal(node)
	res, _ := apiPOST("node", "", params)
	n := NodeList{}
	return n.FromJSONBody(res)
}

// GetNode - GET /api/node/:id
func GetNode(id string) Node {
	res, _ := apiGET("node", id)
	n := NodeList{}
	return n.FromJSONBody(res)
}

// UpdateNode - PATCH /api/node/:id
func UpdateNode(id string, node Node) Node {
	params, _ := json.Marshal(node)
	res, _ := apiPATCH("node", id, params)
	n := NodeList{}
	return n.FromJSONBody(res)
}

// DeleteNode - DELETE /api/node/:id
func DeleteNode(id string) error {
	_, err := apiDELETE("node", id)
	if err != nil {
		return err
	}
	return nil
}

// CreateAccount - POST /api/account
func CreateAccount(account Account) Account {
	params, _ := json.Marshal(account)
	res, _ := apiPOST("account", "", params)
	a := AccountList{}
	return a.FromJSONBody(res)
}

// GetAccount - GET /api/account
func GetAccount() Account {
	res, _ := apiGET("account", "")
	a := AccountList{}
	return a.FromJSONBody(res)
}

// UpdateAccount - PATCH /api/account
func UpdateAccount(account Account) Account {
	params, _ := json.Marshal(account)
	res, _ := apiPATCH("account", "", params)
	a := AccountList{}
	return a.FromJSONBody(res)
}

// DeleteAccount - DELETE /api/account
func DeleteAccount() error {
	_, err := apiDELETE("account", "")
	if err != nil {
		return err
	}
	return nil
}

func apiGET(types string, paths string) (*http.Response, error) {
	authToken, _ := readAuthToken()

	var path string
	if paths != "" {
		path = strings.Join([]string{types, paths}, "/")
	} else {
		path = types
	}
	url := fmt.Sprintf("http://%s/%s", apiHost, path)

	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Auth-Token", authToken)
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func apiPOST(types string, paths string, params []byte) (*http.Response, error) {
	authToken, _ := readAuthToken()

	var path string
	if paths != "" {
		path = strings.Join([]string{types, paths}, "/")
	} else {
		path = types
	}
	url := fmt.Sprintf("http://%s/%s", apiHost, path)

	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(params))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Auth-Token", authToken)
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func apiPATCH(types string, paths string, params []byte) (*http.Response, error) {
	authToken, _ := readAuthToken()

	var path string
	if paths != "" {
		path = strings.Join([]string{types, paths}, "/")
	} else {
		path = types
	}
	url := fmt.Sprintf("http://%s/%s", apiHost, path)

	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(params))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Auth-Token", authToken)
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func apiDELETE(types string, paths string) (*http.Response, error) {
	authToken, _ := readAuthToken()

	var path string
	if paths != "" {
		path = strings.Join([]string{types, paths}, "/")
	} else {
		path = types
	}
	url := fmt.Sprintf("http://%s/%s", apiHost, path)

	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Auth-Token", authToken)
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}
