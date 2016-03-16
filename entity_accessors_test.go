package restful

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type keyvalue struct {
	readCalled  bool
	writeCalled bool
}

func (kv *keyvalue) Read(req *Request, v interface{}) error {
	//t := reflect.TypeOf(v)
	//rv := reflect.ValueOf(v)
	kv.readCalled = true
	return nil
}

func (kv *keyvalue) Write(resp *Response, status int, v interface{}) error {
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

// go test -v -test.run TestKeyValueEncoding ...restful
func TestKeyValueEncoding(t *testing.T) {
	type Book struct {
		Title         string
		Author        string
		PublishedYear int
	}

	kv := new(keyvalue)
	kv2 := new(keyvalue)

	RegisterEntityAccessor("application/vnd.my.company+v1+kv", kv)

	// registering entity accessor for v2
	RegisterEntityAccessor("application/vnd.my.company+v2+kv", kv2)

	b := Book{"Singing for Dummies", "john doe", 2015}

	// Write
	httpWriter := httptest.NewRecorder()

	// SWITCH VERSION
	//								Accept									Produces
	resp := Response{httpWriter, "application/vnd.my.company+v2+kv,*/*;q=0.8", []string{"application/vnd.my.company+v2+kv"}, 0, 0, true, nil}
	// SWITCH VERSION

	resp.WriteEntity(b)


	t.Log(string(httpWriter.Body.Bytes()))

	if !kv.writeCalled {
		t.Error("Write never called - v1")
	}
	if !kv2.writeCalled {
		t.Error("Write never called - v2")
	}

	// Read
	bodyReader := bytes.NewReader(httpWriter.Body.Bytes())
	httpRequest, _ := http.NewRequest("GET", "/test", bodyReader)

	// SWITCH VERSION
	httpRequest.Header.Set("Content-Type", "application/vnd.my.company+v2+kv; charset=UTF-8")
	// SWITCH VERSION

	request := NewRequest(httpRequest)
	var bb Book
	request.ReadEntity(&bb)

	if !kv.readCalled {
		t.Error("Read never called - v1")
	}
	if !kv2.readCalled {
		t.Error("Read never called - v2")
	}
}
