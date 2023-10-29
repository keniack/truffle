package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func IncomingHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if Debug {
			DebugLog.Printf("incoming started")
		}
		target := r.Header.Get("x-target")
		contentLength := r.Header.Get("Content-Length")
		buf := bytes.Buffer{}
		_, _ = io.Copy(&buf, r.Body)
		if Debug {
			DebugLog.Printf("incoming target %s", target)
			DebugLog.Printf("incoming content length %s", contentLength)
		}

		var podIP = GetPodIpForName(target)
		url := fmt.Sprintf("http://%s:%s", podIP, "8080")
		if Debug {
			DebugLog.Printf("url %s", url)
		}
		proxyReq, _ := http.NewRequest(r.Method, url, bytes.NewReader(buf.Bytes()))
		proxyReq.Header.Add("Content-Type", "application/json")
		proxyReq.Header.Add("x-target", target)
		client := GetHttpClient()
		for {
			resp, err := client.Do(proxyReq)
			if err != nil {
				continue
			}
			_, err = io.Copy(w, resp.Body)

			break
		}
		if Debug {
			DebugLog.Printf("incoming finished. response OK")
		}

	}
}

func OutgoingHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if Debug {
			DebugLog.Printf("outgoing started")
		}
		target := r.Header.Get("x-target")
		contentLength := r.Header.Get("Content-Length")
		//body, err := io.ReadAll(r.Body)
		//r.Body = io.NopCloser(bytes.NewReader(body))
		buf := bytes.Buffer{}
		_, _ = io.Copy(&buf, r.Body)
		if Debug {
			DebugLog.Printf("outgoing target %s", target)
			DebugLog.Printf("outgoing content length %s", contentLength)
		}
		var nodeIp = GetNodeIpForName(target)
		url := fmt.Sprintf("http://%s:%s", nodeIp, "8888/hello")

		if Debug {
			DebugLog.Printf("url %s", url)
		}
		//proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
		proxyReq, _ := http.NewRequest(r.Method, url, bytes.NewReader(buf.Bytes()))
		proxyReq.Header.Add("Content-Type", "application/json")
		proxyReq.Header.Add("x-target", target)
		client := GetHttpClient()
		resp, err := client.Do(proxyReq)
		//respBody, _ := io.ReadAll(resp.Body)
		//log.Printf("outgoing finished. response OK %s", string(respBody))
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

	}
}

func ProxyIncomingHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if Debug {
			log.Println("Reverse incoming proxy started")
		}
		target := r.Header.Get("x-target")
		if Debug {
			contentLength := r.Header.Get("Content-Length")
			log.Printf("incoming target %s", target)
			log.Printf("incoming content length %s", contentLength)
		}

		var podIP = GetPodIpForName(target)
		targetIpPort := fmt.Sprintf("%s:%s", podIP, "8080")

		//targetIpPort := fmt.Sprintf("%s:%s", podIP, "80")
		newTargetUrl := fmt.Sprintf("http://%s", targetIpPort)
		targetUrl, _ := url.Parse(newTargetUrl)
		newRequestUrl, _ := url.Parse("")
		r.URL = newRequestUrl
		p.Director = func(req *http.Request) {
			RewriteRequestURL(req, targetUrl)
		}
		if Debug {
			log.Printf("start test incoming target %s", newTargetUrl)
		}
		for {
			if IsTCPWorking(targetIpPort) {
				break
			}
			continue
		}
		if Debug {
			log.Println("Forward incoming finished")
		}
		p.ServeHTTP(w, r)
	}
}

func ProxyOutgoingHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if Debug {
			log.Println("Reverse outgoing proxy started")
		}
		target := r.Header.Get("x-target")
		if Debug {
			contentLength := r.Header.Get("Content-Length")
			log.Printf("outgoing target %s", target)
			log.Printf("outgoing content length %s", contentLength)
		}
		var nodeIp = GetNodeIpForName(target)
		newTargetUrl := fmt.Sprintf("http://%s:%s", nodeIp, "8888")
		if Debug {
			log.Printf("outgoing target %s", newTargetUrl)
		}

		//newTargetUrl := fmt.Sprintf("http://%s:%s", "192.168.0.207", "8888")
		targetUrl, _ := url.Parse(newTargetUrl)
		newRequestUrl, _ := url.Parse("/hello")
		r.URL = newRequestUrl
		p.Director = func(req *http.Request) {
			RewriteRequestURL(req, targetUrl)
		}
		if Debug {
			DebugLog.Println("outgoing finished")
		}
		p.ServeHTTP(w, r)
	}
}
