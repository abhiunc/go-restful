package main

import (
	"fmt"
	"io"
	"github.com/emicklei/go-restful"
	"net/http"
	"reflect"
)

// simple struct for registering new entity accessor
type keyvalue struct {
	readCalled bool
	writeCalled bool
}

// generic read function
func (kv *keyvalue) Read(req *restful.Request, v interface{}) error {
	kv.readCalled = true
	return nil
}

// generic write function
func (kv *keyvalue) Write(resp *restful.Response, status int, v interface{}) error {
	t := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	for ix := 0; ix < t.NumField(); ix++ {
		sf := t.Field(ix)
		io.WriteString(resp, sf.Name)
		io.WriteString(resp, "=")
		io.WriteString(resp, fmt.Sprintf("%v\n", rv.Field(ix).Interface()))
	}
	kv.writeCalled = true
	return nil
}

func main() {
	kv := new(keyvalue)

	// create new entity accessor for application/test+v1 since it is not standard
	restful.RegisterEntityAccessor("application/test+v1", kv)

	// spin up a new webservice
	ws := new(restful.WebService)

	// create path to /testing
	ws.Path("/testing").
		Consumes(restful.MIME_JSON). 	// standard json input
		Produces(restful.MIME_JSON)		// standard json output

	ws.Route(ws.GET("/").To(hello).			// route GET at root of /testing to hello function
		Consumes("application/test+v1").	// consumes our custom accept header
		Produces("application/test+v1"))	// produces our custom accept header

	restful.Add(ws)
	http.ListenAndServe(":8080", nil)
}

// checks the Accept header of request
func hello(req *restful.Request, resp *restful.Response) {
	// fmt.Println(req.Request.Header.Get(restful.HEADER_Accept))

	// print VERSION 1 if header == application/test+v2
	if(req.Request.Header.Get(restful.HEADER_Accept) == "application/test+v1") {
		fmt.Println("VERSION 1")
	} else { // print DEFAULT if header == "" or "*/*"
		fmt.Println("DEFAULT")
	}

	// any other header will be unaccepted!
}