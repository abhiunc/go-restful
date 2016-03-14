package main

import (
	"io"
	"net/http"

	"github.com/abhiunc/go-restful"
)

// This example shows the minimal code needed to get a restful.WebService working.
//
// GET http://localhost:8080/hello

func main() {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/hello").To(AuthAspect(hello))) //wrap aspect around RouteFunction
	restful.Add(ws)
	http.ListenAndServe(":8000", nil)
}

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world")
}
