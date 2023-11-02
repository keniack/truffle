package common

import (
	"log"
)

var RedisIP string
var RedisPwd string
var AwsAccessKey string
var AwsSecretKey string
var ComMode string
var IncomingPodPort string
var Debug, Trace bool
var DebugLog, TraceLog *log.Logger
