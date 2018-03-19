package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	uuid "github.com/satori/go.uuid"
)

type Request struct {
	Filename string `json:"filename"`
}

type Response struct {
	URL string `json:"url"`
	Key string `json:"key"`
	Ok  bool   `json:"ok"`
}

func generateFilename(filename string) string {
	ext := filepath.Ext(filename)
	u := uuid.Must(uuid.NewV4())
	secureFilename := fmt.Sprintf("%s%s", u, ext)
	return secureFilename
}

func Handler(request Request) (Response, error) {
	filename := generateFilename(request.Filename)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:        aws.String("recognition-images"),
		Key:           aws.String(filename),
		ContentLength: aws.Int64(0),
	})
	url, _, err := req.PresignRequest(15 * time.Minute)

	if err != nil {
		fmt.Println("error presigning request", err)
		return Response{Ok: false}, err
	}

	return Response{
		URL: url,
		Key: filename,
		Ok:  true,
	}, nil
}

func main() {
	// filename := os.Args[1]
	// resp, _ := Handler(Request{Filename: filename})
	// fmt.Println(resp.Key)
	// fmt.Println(resp.URL)

	lambda.Start(Handler)
}
