package main

import (
	"bytes"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Request is a input request structure
type Request struct {
	SourceBucket string `json:"sourcebucket"`
	DestBucket   string `json:"destinationbucket"`
}

// HandlerRequest handles incoming requests
func HandlerRequest(req Request) (out string, err error) {

	log.Print("Starting")
	//process sync/async
	process(req)

	log.Print("Done")

	return
}

func main() {
	lambda.Start(HandlerRequest)

	// start := time.Now()
	// req := Request{SourceBucket: "testlambdaimages", DestBucket: "testlambdaimages-small"}
	// HandlerRequest(req)
	//
	// log.Printf("Execution Time: %s", time.Since(start))

}

var svc *s3.S3

func process(req Request) {
	//svc := s3.New(session.New())
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("S3_REGION"))},
	)
	svc = s3.New(awsSession)
	// Get the list of objects need to process
	var s3Objects []*s3.Object
	if s3Objects, err = getS3BucketOjects(svc, req.SourceBucket); err == nil {
		// Process objects in bucket.
		err = processObjects(svc, req, s3Objects)
	}
	if err != nil {
		log.Printf("Error- %v", err)
	}
}

func getS3BucketOjects(svc *s3.S3, bucketName string) (s3Objects []*s3.Object, err error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}
	var result *s3.ListObjectsV2Output
	if result, err = svc.ListObjectsV2(input); err == nil {
		s3Objects = result.Contents
	}
	return
}

var wg sync.WaitGroup

type cReq struct {
	obj *s3.GetObjectOutput
	key string
	req Request
}

func processObjects(svc *s3.S3, req Request, s3Objects []*s3.Object) (err error) {

	var ch = make(chan cReq, len(s3Objects))
	wg.Add(len(s3Objects))

	for _, f := range s3Objects {
		go func(ch chan cReq, f *s3.Object) {
			log.Print("Downloading- ", *f.Key)
			r, out := svc.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(req.SourceBucket),
				Key:    aws.String(*f.Key),
			})
			err := r.Send()
			if err != nil {
				log.Fatal("Error while downloading", *f.Key, err)
			}
			hh := cReq{obj: out, key: *f.Key}

			ch <- hh
		}(ch, f)

	}
	for i := 0; i < len(s3Objects); i++ {
		select {
		case m := <-ch:
			go processAync(m.obj, m.key, req.DestBucket)
			// case <-time.After(3 * time.Second):
			// 	fmt.Println("timeout 2")
		}
	}

	wg.Wait()
	//close(ch)
	return
}

func processAync(obj *s3.GetObjectOutput, key, destBucket string) {

	img, err := jpeg.Decode(obj.Body)
	if err != nil {
		log.Fatal(err)
	}

	imgOut, _ := Resize(img, 30, 30)

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, imgOut, nil)

	go func(buf bytes.Buffer, key, destBucket string) {
		defer wg.Done()
		reqUpload, _ := svc.PutObjectRequest(&s3.PutObjectInput{
			Bucket:      aws.String(destBucket),
			Key:         aws.String(key),
			Body:        bytes.NewReader(buf.Bytes()),
			ContentType: aws.String(http.DetectContentType(buf.Bytes())),
		})
		err = reqUpload.Send()
		if err != nil {
			log.Fatal("Error while uploading", key, err)
		} else {
			log.Print("Completed- ", key)
		}
	}(*buf, key, destBucket)

}
