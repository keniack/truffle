package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"polaris/truffle/pkg/common"
	"polaris/truffle/pkg/utils"
)

func ProxyIncomingHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if common.Debug {
			log.Println("Reverse incoming proxy started")
		}
		target := r.Header.Get("x-target")
		if common.Debug {
			contentLength := r.Header.Get("Content-Length")
			log.Printf("incoming target %s", target)
			log.Printf("incoming content length %s", contentLength)
		}

		//var podIP = GetPodIpForName(target)
		var podIP = "GetPodIpForName(target)"
		targetIpPort := fmt.Sprintf("%s:%s", podIP, "80")

		//targetIpPort := fmt.Sprintf("%s:%s", podIP, "80")
		newTargetUrl := fmt.Sprintf("http://%s", targetIpPort)
		targetUrl, _ := url.Parse(newTargetUrl)
		newRequestUrl, _ := url.Parse("")
		r.URL = newRequestUrl
		p.Director = func(req *http.Request) {
			utils.RewriteRequestURL(req, targetUrl)
		}
		if common.Debug {
			log.Printf("start testing connection to incoming target %s", newTargetUrl)
		}
		for {
			if IsTCPWorking(targetIpPort) {
				break
			}
			continue
		}
		if common.Debug {
			log.Println("Forward incoming finished")
		}
		p.ServeHTTP(w, r)
	}
}

func ProxyOutgoingHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if common.Debug {
			log.Println("Reverse outgoing proxy started")
		}
		target := r.Header.Get("x-target")
		if common.Debug {
			contentLength := r.Header.Get("Content-Length")
			log.Printf("outgoing target %s", target)
			log.Printf("outgoing content length %s", contentLength)
		}
		var nodeIp = "127.0.0.1"
		//var nodeIp = GetNodeIpForName(target)
		newTargetUrl := fmt.Sprintf("http://%s:%s", nodeIp, "8888")
		if common.Debug {
			log.Printf("outgoing target %s", newTargetUrl)
		}

		//newTargetUrl := fmt.Sprintf("http://%s:%s", "192.168.0.207", "8888")
		targetUrl, _ := url.Parse(newTargetUrl)
		newRequestUrl, _ := url.Parse("/hello")
		r.URL = newRequestUrl
		p.Director = func(req *http.Request) {
			utils.RewriteRequestURL(req, targetUrl)
		}
		if common.Debug {
			common.DebugLog.Println("outgoing finished")
		}
		p.ServeHTTP(w, r)
	}
}
