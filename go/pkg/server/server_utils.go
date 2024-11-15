package server

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func NewProxy(targetUrl string) (*httputil.ReverseProxy, error) {
	target, _ := url.Parse(targetUrl)
	proxy := httputil.NewSingleHostReverseProxy(target)
	return proxy, nil
}

func GetHttpClient() http.Client {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	transport := &http.Transport{
		DialContext:        dialer.DialContext,
		DisableKeepAlives:  true,
		DisableCompression: true,
	}
	client := http.Client{
		Transport: transport,
	}
	return client
}

func IsTCPWorking(url string) bool {
	conn, err := net.Dial("tcp", url)
	if conn != nil {
		err := conn.Close()
		if err != nil {
			return false
		} // close if problem
	}
	if err != nil {
		return false
	}
	return true
}
