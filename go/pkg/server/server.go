package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"polaris/truffle/pkg/client"
	"polaris/truffle/pkg/common"
	"polaris/truffle/pkg/watcher"
	"strconv"
	"sync"
)

func OutgoingHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if common.Debug {
			common.DebugLog.Printf("outgoing started")
		}
		target := r.Header.Get("x-target")
		contentLengthStr := r.Header.Get("Content-Length")
		contentLength, _ := strconv.Atoi(contentLengthStr)
		if common.Debug {
			common.DebugLog.Printf("outgoing target %s", target)
			common.DebugLog.Printf("Comm Mode %s outgoing content length %s", common.ComMode, contentLengthStr)
		}
		var nodeIpChannel = make(chan string, 1)
		var contentChannel = make(chan bytes.Buffer, 1)

		var wg sync.WaitGroup
		wg.Add(2)
		//buf := *bytes.NewBuffer(make([]byte, 0, contentLength))
		go client.SetContentOutgoing(r.Body, contentChannel, contentLength, &wg)
		go watcher.GetNodeIpForName(target, nodeIpChannel, &wg)
		buf := <-contentChannel
		nodeIp := <-nodeIpChannel
		wg.Wait()
		close(nodeIpChannel)
		close(contentChannel)

		url := fmt.Sprintf("http://%s:%s", nodeIp, "8888/hello")

		if common.Debug {
			common.DebugLog.Printf("url %s", url)
		}
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

func IncomingHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if common.Debug {
			common.DebugLog.Printf("incoming started")
		}
		target := r.Header.Get("x-target")
		var podIP string
		var byteContent []byte

		var wg sync.WaitGroup
		wg.Add(2)
		var contentChannel = make(chan []byte, 1)
		var podIpChannel = make(chan string, 1)
		go client.GetContentIncoming(r.Body, contentChannel, &wg)
		go watcher.GetPodIpForName(target, podIpChannel, &wg)
		byteContent = <-contentChannel
		podIP = <-podIpChannel
		wg.Wait()
		close(podIpChannel)
		close(contentChannel)
		if common.Debug {
			common.DebugLog.Printf("incoming target %s", target)
			common.DebugLog.Printf("Comm Mode %s incoming content length %d", common.ComMode, len(byteContent))
		}

		url := fmt.Sprintf("http://%s:%s", podIP, common.IncomingPodPort)
		if common.Debug {
			common.DebugLog.Printf("url %s", url)

		}
		proxyReq, _ := http.NewRequest(r.Method, url, bytes.NewReader(byteContent))
		proxyReq.Header.Add("Content-Type", "application/json")
		proxyReq.Header.Add("x-target", target)
		client := GetHttpClient()
		for {
			resp, err := client.Do(proxyReq)
			if err != nil {
				continue
			}
			_, err = io.Copy(w, resp.Body)
			defer resp.Body.Close()
			break
		}
		if common.Debug {
			common.DebugLog.Printf("incoming finished. response OK")
		}

	}
}
