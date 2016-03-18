package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

func renderStuff(w http.ResponseWriter, stuff interface{}) {
	data, _ := json.Marshal(stuff)
	jsonString := string(data) + "\n"

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Fprint(w, jsonString)
}

func renderResponse(w http.ResponseWriter, status int, message string, content interface{}) {
	if content == nil {
		content = map[string]interface{}{}
	}
	renderStuff(w, jsonResponse{Status: status, Message: message, Content: content})
}

func renderError(w http.ResponseWriter, err error) {
	renderResponse(w, 500, err.Error(), nil)
}

func renderContent(w http.ResponseWriter, content interface{}) {
	renderResponse(w, 200, "OK", content)
}

func retrievePostJSON(r *http.Request, v interface{}) {
	jsonIn := r.PostFormValue("json")

	json.Unmarshal([]byte(jsonIn), v)
}

type serverActions struct {
	BOINCAddr  string
	Connection *BOINCConn
	Client     *Client
}

func (s *serverActions) DoAuth(w http.ResponseWriter, r *http.Request) error {
	var postData PostRequest
	retrievePostJSON(r, &postData)

	boincAddrP := postData.BOINCAddr
	rpcKeyP := postData.Password
	switch {
	case rpcKeyP == nil:
		return errors.New("No RPC key specified")
	case boincAddrP == nil:
		return errors.New("BOINC host to query not specified")
	default:
		var URL, URLErr = url.Parse(*boincAddrP)
		if URLErr != nil {
			return errors.New("Invalid URL specified")
		}
		return s.Client.Connect(*URL, &(*rpcKeyP))
	}
}

func (s *serverActions) AuthExec(w http.ResponseWriter, r *http.Request, cb func()) {
	err := s.DoAuth(w, r)
	if err == nil {
		if cb != nil {
			cb()
		}
		s.Client.Disconnect()
	} else {
		renderResponse(w, 500, err.Error(), nil)
	}
}

func (s *serverActions) TryAuth(w http.ResponseWriter, r *http.Request) {
	s.AuthExec(w, r, func() { renderResponse(w, 200, "OK", nil) })
}

func (s *serverActions) GetStuff(w http.ResponseWriter, r *http.Request, fn func() (interface{}, error), k string) {
	s.AuthExec(w, r, func() {
		stuff, err := fn()
		if err != nil {
			renderError(w, err)
		} else {
			renderContent(w, map[string]interface{}{k: stuff})
		}
	})
}

func (s *serverActions) GetProjects(w http.ResponseWriter, r *http.Request) {
	s.GetStuff(w, r, func() (interface{}, error) { return s.Client.GetProjects() }, "projects")
}

func (s *serverActions) GetMessages(w http.ResponseWriter, r *http.Request) {
	s.GetStuff(w, r, func() (interface{}, error) { return s.Client.GetMessages() }, "messages")
}

func (s *serverActions) GetAcctMgrInfo(w http.ResponseWriter, r *http.Request) {
	s.GetStuff(w, r, func() (interface{}, error) { return s.Client.GetAcctMgrInfo() }, "acct_mgr_info")
}

func (s *serverActions) PollAcctMgrRPC(w http.ResponseWriter, r *http.Request) {
	s.GetStuff(w, r, func() (interface{}, error) { return s.Client.PollAcctMgrRPC() }, "acct_mgr_rpc")
}

func startActionInstance(ClientInstance *Client) serverActions {
	return serverActions{Client: ClientInstance}
}
