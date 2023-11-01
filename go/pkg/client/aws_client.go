package client

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"polaris/truffle/pkg/common"
)

const (
	S3_REGION = "eu-central-1"
	S3_BUCKET = "truffle-video-stream"
)

type S3Handler struct {
	Session *session.Session
	Bucket  string
}

func (h S3Handler) GetFile(key string) bytes.Buffer {

	output, _ := s3.New(h.Session).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(h.Bucket),
		Key:    aws.String(key),
	})
	buf := bytes.Buffer{}
	_, _ = io.Copy(&buf, output.Body)
	return buf
}

func (h S3Handler) UploadFile(key string, content []byte) {

	_, _ = s3.New(h.Session).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(h.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(content),
	})
}

func GetValueS3AWS(key string) []byte {
	file := getS3Handler().GetFile(key)
	return file.Bytes()
}

func getS3Handler() S3Handler {
	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(S3_REGION),
			Credentials: credentials.NewStaticCredentials(common.AwsAccessKey, common.AwsSecretKey, ""),
		})
	if err != nil {
		fmt.Printf("session.NewSession, err: %v\n", err)
	}

	handler := S3Handler{
		Session: sess,
		Bucket:  S3_BUCKET,
	}
	return handler
}

func SetValueS3AWS(key string, content []byte) {
	getS3Handler().UploadFile(key, content)
}
