package client

import (
	"bytes"
	"io"
	"polaris/truffle/pkg/common"
	"sync"
)

const FILENAME = "truffleTempFile"

func SetContentOutgoing(contentBody io.Reader, contentChannel chan<- bytes.Buffer, wg *sync.WaitGroup) {
	defer wg.Done()
	if common.Debug {
		common.DebugLog.Printf("Started copying outgoing")
	}
	buf := bytes.Buffer{}
	_, _ = io.Copy(&buf, contentBody)
	switch {
	case common.ComMode == "S3":
		SetValueS3AWS(FILENAME, buf.Bytes())
		contentChannel <- bytes.Buffer{}
	case common.ComMode == "KVS":
		SetKeyKVS(FILENAME, buf.Bytes())
		contentChannel <- bytes.Buffer{}
	default:
		contentChannel <- buf
	}
	if common.Debug {
		common.DebugLog.Printf("finished copying outgoing")
	}
}

func GetContentIncoming(contentBody io.Reader, contentChannel chan<- []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	if common.Debug {
		common.DebugLog.Printf("Started copying incoming")
	}
	buf := bytes.Buffer{}
	switch {
	case common.ComMode == "S3":
		contentChannel <- GetValueS3AWS(FILENAME)
	case common.ComMode == "KVS":
		contentChannel <- GetKeyKVS(FILENAME)
	default:
		_, _ = io.Copy(&buf, contentBody)
		contentChannel <- buf.Bytes()
	}
	if common.Debug {
		common.DebugLog.Printf("finished copying incoming")
	}
}
