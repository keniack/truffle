package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"polaris/truffle/api"
)

const IncomingPath = "/hello"
const OutgoingPath = "/outgoing"

var (
	ns, label, field             string
	incomingProxy, outgoingProxy bool
)

func init() {

	flag.StringVar(&ns, "namespace", "", "namespace")
	flag.StringVar(&label, "l", "", "Label selector")
	flag.StringVar(&field, "f", "", "Field selector")
	flag.BoolVar(&api.Debug, "debug", false, "Debug log level")
	flag.BoolVar(&api.Trace, "trace", false, "Trace log level")
	flag.BoolVar(&incomingProxy, "incoming-proxy", false, "Use reverse proxy for incoming requests")
	flag.BoolVar(&outgoingProxy, "outgoing-proxy", false, "Use reverse proxy for outgoing requests")
	flag.Parse()
	if api.Trace {
		api.Debug = true
	}

	logWriter := new(api.LogWriter)
	api.InfoLog = log.New(logWriter, "[INFO] ", 0)
	api.DebugLog = log.New(logWriter, "[DEBUG] ", 0)
	api.TraceLog = log.New(logWriter, "[TRACE] ", 0)
	api.ErrorLog = log.New(logWriter, "[ERROR] ", 0)

	log.SetFlags(0)
	log.SetOutput(logWriter)

}

// kubectl get pods
func main() {
	if api.Debug {
		api.DebugLog.Printf("This is something the label set %s", label)
	}

	r := mux.NewRouter()
	// Returns a proxy for the target url.
	proxy, err := api.NewProxy("http://localhost:8080")
	if err != nil {
		panic(err)
	}
	if r.HandleFunc(IncomingPath, api.IncomingHandler()); incomingProxy {
		r.HandleFunc(IncomingPath, api.ProxyIncomingHandler(proxy))
	}
	if r.HandleFunc(OutgoingPath, api.OutgoingHandler()); outgoingProxy {
		r.HandleFunc(OutgoingPath, api.ProxyOutgoingHandler(proxy))
	}
	err = http.ListenAndServe(":8888", r)
	if err != nil {
		return
	}

}
