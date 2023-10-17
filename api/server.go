package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func StartServer() {
	log.Printf("start webserver\n")

	incoming()
	outgoing()
	log.Printf("Starting server at port 8888\n")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}

}

func incoming() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("incoming started")
		reqBody, _ := io.ReadAll(r.Body)
		var message Message
		s, _ := strconv.Unquote(string(reqBody))
		err := json.Unmarshal([]byte(s), &message)
		if err != nil {
			log.Printf("Error parsing request body %s", err)
		}
		var start = time.Now()
		for !ExistsAndHasPodIp(PodsMap, message.Target) {
			//log.Printf("incoming %s pod not found wait 10ms", message.Target)
			time.Sleep(10 * time.Millisecond)
			if start.Add(time.Second * 10).Before(time.Now()) {
				break
			}
		}
		pod, _ := PodsMap.Load(message.Target)
		resBody := client(reqBody, pod.PodIp, ":8080")
		_, err = fmt.Fprintf(w, "Hello "+resBody)
		log.Printf("incoming finished. response %s", resBody)
		if err != nil {
			return
		}
	})
}

func outgoing() {
	http.HandleFunc("/outgoing", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("outgoing started")
		reqBody, _ := io.ReadAll(r.Body)
		var message Message
		s, _ := strconv.Unquote(string(reqBody))
		err := json.Unmarshal([]byte(s), &message)
		if err != nil {
			log.Printf("Error parsing request body %s", err)
		}
		var start = time.Now()
		for !ExistsInMap(PodsMap, message.Target) {
			//log.Printf("outoing %s pod not found wait 10ms", message.Target)
			time.Sleep(10 * time.Millisecond)
			if start.Add(time.Second * 10).Before(time.Now()) {
				break
			}
		}
		pod, _ := PodsMap.Load(message.Target)
		resBody := client(reqBody, pod.NodeIP, ":8888/hello")
		_, err = fmt.Fprintf(w, "Hello "+resBody)
		log.Printf("outgoing finished. response %s", resBody)
		if err != nil {
			return
		}
	})
}

func client(data []byte, targetIp string, port string) string {
	client := http.Client{Timeout: time.Duration(10) * time.Second}
	req, _ := http.NewRequest("POST", "http://"+targetIp+port, bytes.NewBuffer(data))
	req.Header.Add("Content-Type", "application/json")
	log.Printf("started client")
	for {
		resp, err := client.Do(req)
		if err != nil {
			//fmt.Printf("error sending the first time: %v\n", err)
			continue
		}

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		resBody, err := io.ReadAll(resp.Body)
		return string(resBody)
	}
}
