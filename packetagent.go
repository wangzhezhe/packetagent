package main

import (
	"flag"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
	"packetagent/app"
)

// This example shows the minimal code needed to get a restful.WebService working.
//
// GET http://localhost:8080/hello
var (
	Listenport  = "9999"
	Defaulttime = 60
)

func main() {
	flag.Parse()
	wsContainer := restful.NewContainer()
	if glog.V(1) {
		glog.Info("start listen packet")
	}
	app.Register(wsContainer)

	server := &http.Server{Addr: ":" + Listenport, Handler: wsContainer}
	server.ListenAndServe()
}
