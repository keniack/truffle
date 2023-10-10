package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

		reqBody, _ := io.ReadAll(r.Body)

		var message Message
		err := json.Unmarshal(reqBody, &message)
		if err != nil {
			log.Printf("Error parsing request body %s", err)
		}
		var start = time.Now()
		for !ExistsAndHasPodIp(PodsMap, message.Target) {
			log.Print("pod not found wait 10ms")
			time.Sleep(10 * time.Millisecond)
			if start.Add(time.Second * 10).Before(time.Now()) {
				break
			}
		}
		pod := PodsMap[message.Target]
		resBody := client(reqBody, pod.PodIp, ":8080")
		_, err = fmt.Fprintf(w, "Hello "+resBody)
		if err != nil {
			return
		}
	})
}

func outgoing() {
	http.HandleFunc("/outgoing", func(w http.ResponseWriter, r *http.Request) {

		reqBody, _ := io.ReadAll(r.Body)
		var message Message
		err := json.Unmarshal(reqBody, &message)
		if err != nil {
			log.Printf("Error parsing request body %s", err)
		}
		var start = time.Now()
		for !ExistsInMap(PodsMap, message.Target) {
			log.Print("pod not found wait 10ms")
			time.Sleep(10 * time.Millisecond)
			if start.Add(time.Second * 10).Before(time.Now()) {
				break
			}
		}
		pod := PodsMap[message.Target]
		resBody := client(reqBody, pod.NodeIP, ":8888/hello")
		_, err = fmt.Fprintf(w, "Hello "+resBody)
		if err != nil {
			return
		}
	})
}

func client(data []byte, targetIp string, port string) string {
	c := http.Client{Timeout: time.Duration(1) * time.Hour}
	resp, err := c.Post("http://"+targetIp+port, "text/plain", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error %s", err)
		return ""
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	resBody, err := io.ReadAll(resp.Body)
	return string(resBody)
}
