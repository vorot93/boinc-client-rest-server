package main

import (
	"flag"
	"net/http"
)

func main() {
	var exitChan = make(chan struct{})
	var waitForExit = func() { <-exitChan }
	var startExitSequence = func() { close(exitChan) }

	var shimAddr = flag.String("shimaddr", ":15678", "Address of shim server.")
	flag.Parse()
	core := MakeClient()
	actions := startActionInstance(core)

	const APIVer = "0.1"
	var APIPrefix = "/" + APIVer

	sMux := http.NewServeMux()
	var handlerMap = map[string]func(http.ResponseWriter, *http.Request){}
	handlerMap["quit"] = func(w http.ResponseWriter, r *http.Request) { startExitSequence() }
	handlerMap["tryauth"] = actions.TryAuth
	handlerMap["projects"] = actions.GetProjects
	handlerMap["messages"] = actions.GetMessages
	handlerMap["acctmgr"] = actions.PollAcctMgrRPC
	handlerMap["acctmgr/info"] = actions.GetAcctMgrInfo
	handlerMap["entrypoints"] = func(w http.ResponseWriter, r *http.Request) {
		var entrypoints = make([]string, 0, len(handlerMap))
		for k := range handlerMap {
			entrypoints = append(entrypoints, k)
		}

		renderContent(w, map[string][]string{"entrypoints": entrypoints})
	}

	for k, v := range handlerMap {
		sMux.HandleFunc(APIPrefix+"/"+k, v)
	}

	server := &http.Server{
		Addr:    *shimAddr,
		Handler: sMux,
	}

	go server.ListenAndServe()

	waitForExit()
}
