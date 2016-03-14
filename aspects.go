package restful

import (
	"log"
	"net/http"
	"net/http/pprof"
	"strconv"
	"strings"

	"github.com/abhiunc/go-restful"
)

//Creates Logger before RouteFunction is executed via wrapping
func LoggingAspect(fn restful.RouteFunction) restful.RouteFunction {
	return func(r *restful.Request, w *restful.Response) {
		log.Printf("[WebserviceLogging] {%s}, \n%s  \n%s  \n%s  \n%s  \n%s\n",
			strings.Split(r.Request.RemoteAddr, ":")[0],
			"[Method]: "+r.Request.Method,
			"[URL]: "+r.Request.URL.String(),
			"[Protocol]: "+r.Request.Proto,
			"[Status]: "+strconv.Itoa(w.StatusCode()),
			"[Length]: "+strconv.Itoa(w.ContentLength()),
		)
		fn(r, w)
	}
}

//Creates Authentication at RouteFunction level.
func AuthAspect(fn restful.RouteFunction) restful.RouteFunction {
	return func(r *restful.Request, w *restful.Response) {
		token := r.HeaderParameter("X-AUTH-TOKEN")
		if token == "" {
			http.Error(w, "missing auth token", http.StatusUnauthorized)
			return
		}
		//Write own authentication mechanism here.
		fn(r, w)
	}
}

//runs Index profiling on wrapped RouteFunction
func IndexHandler(fn restful.RouteFunction) restful.RouteFunction {
	return func(r *restful.Request, w *restful.Response) {
		pprof.Index(w.ResponseWriter, r.Request)
		fn(r, w)
	}
}

//Trace profiler on wrapper RouteFunction
func TraceHandler(fn restful.RouteFunction) restful.RouteFunction {
	return func(r *restful.Request, w *restful.Response) {
		pprof.Trace(w.ResponseWriter, r.Request)
		fn(r, w)
	}
}

//General Profiling on RouteFunction
func ProfileHandler(fn restful.RouteFunction) restful.RouteFunction {
	return func(r *restful.Request, w *restful.Response) {
		pprof.Profile(w.ResponseWriter, r.Request)
		fn(r, w)
	}
}

//Symbol Profiling on RouteFunction
func SymbolHandler(fn restful.RouteFunction) restful.RouteFunction {
	return func(r *restful.Request, w *restful.Response) {
		//c := pprof.Symbol(w.ResponseWriter, r.Request)
		log.Printf(" ")
		fn(r, w)
	}
}
