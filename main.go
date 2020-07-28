package main

import (
	"bytes"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Request is a input request structure
type Request struct {
	SourceBucket string `json:"sourcebucket"`
	DestBucket   string `json:"destinationbucket"`
	Timeout      int    `json:"timeout"`
}

// Response is a output reponse structure
type Response struct {
	Message string `json:"message"`
}

// HandlerRequest handles incoming requests
func HandlerRequest(req Request) (out string, err error) {

	log.Print("Starting")
	//svc := s3.New(session.New())
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("S3_REGION"))},
	)
	svc := s3.New(awsSession)
	// Get the list of objects need to process
	var s3Objects []*s3.Object
	if s3Objects, err = getS3BucketOjects(svc, req.SourceBucket); err == nil {
		// Process objects in bucket.
		err = processObjects(svc, req, s3Objects)
	}
	if err != nil {
		log.Printf("Error- %v", err)
	}
	log.Print("Done")

	return
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

func processObjects(svc *s3.S3, req Request, s3Objects []*s3.Object) (err error) {

	for _, f := range s3Objects {
		log.Print("Processing- ", *f.Key)
		r, out := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(req.SourceBucket),
			Key:    aws.String(*f.Key),
		})
		err := r.Send()
		if err != nil {
			log.Fatal("Error while downloading", *f.Key, err)
		}

		img, err := jpeg.Decode(out.Body)
		if err != nil {
			log.Fatal(err)
		}
		imgOut, _ := ResizeImg(img, 30, 30)
		buf := new(bytes.Buffer)
		jpeg.Encode(buf, imgOut, nil)

		r1, _ := svc.PutObjectRequest(&s3.PutObjectInput{
			Bucket:      aws.String(req.DestBucket),
			Key:         aws.String(*f.Key),
			Body:        bytes.NewReader(buf.Bytes()),
			ContentType: aws.String(http.DetectContentType(buf.Bytes())),
		})
		err = r1.Send()
		if err != nil {
			log.Fatal("Error while uploading", *f.Key, err)
		} else {
			log.Print("Completed- ", *f.Key)
		}

	}

	return
}

func main() {
	start := time.Now()
	// req := Request{SourceBucket: "testlambdaimages", DestBucket: "testlambdaimages-small"}
	// _, err := HandlerRequest(req)
	lambda.Start(HandlerRequest)

	log.Printf("Execution Time: %s", time.Since(start))

}
