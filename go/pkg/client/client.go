package client

import (
	"bytes"
	"io"
	"polaris/truffle/pkg/common"
	"sync"
)

const FILENAME = "truffleTempFile"

func SetContentOutgoing(contentBody io.Reader, contentChannel chan<- bytes.Buffer, contentLength int, wg *sync.WaitGroup) {
	defer wg.Done()
	if common.Debug {
		common.DebugLog.Printf("Started copying outgoing")
	}
	buf := *bytes.NewBuffer(make([]byte, 0, contentLength))
	_, _ = io.Copy(&buf, contentBody)
	switch {
	case common.ComMode == "S3":
		if common.Debug {
			common.DebugLog.Printf("start copying in s3")
		}
		SetValueS3AWS(FILENAME, buf.Bytes())
		if common.Debug {
			common.DebugLog.Printf("finished copying in s3")
		}
		contentChannel <- *bytes.NewBuffer(make([]byte, 0, 0))
	case common.ComMode == "KVS":
		if common.Debug {
			common.DebugLog.Printf("start copying in redis")
		}
		SetKeyKVS(FILENAME, buf.Bytes())
		if common.Debug {
			common.DebugLog.Printf("finished copying in redis")
		}
		contentChannel <- *bytes.NewBuffer(make([]byte, 0, 0))
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

	switch {
	case common.ComMode == "S3":
		contentChannel <- GetValueS3AWS(FILENAME)
	case common.ComMode == "KVS":
		contentChannel <- GetKeyKVS(FILENAME)
	default:
		buf := bytes.Buffer{}
		_, _ = io.Copy(&buf, contentBody)
		contentChannel <- buf.Bytes()
	}
	if common.Debug {
		common.DebugLog.Printf("finished copying incoming")
	}
}
