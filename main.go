package main

import (
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"

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
func HandlerRequest(req Request) (out Response, err error) {

	// log.Print(req)
	svc := s3.New(session.New())

	// Get the list of objects need to process
	var s3Objects []*s3.Object
	if s3Objects, err = getS3BucketOjects(svc, req.SourceBucket); err != nil {
		// Process objects in bucket.
		out, err = processObjects(svc, s3Objects)
	}

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

func processObjects(svc *s3.S3, s3Objects []*s3.Object) (outputMsg Response, err error) {
	path := "./images_in/"
	pathOut := "./images_out/"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {

		log.Print("Processing ", f.Name())
		file, err := os.Open(path + f.Name())
		if err != nil {
			log.Fatal(err)
		}

		img, err := jpeg.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		imgOut, _ := ResizeImg(img, 30, 30)

		//this goes to destination bucket
		out, err := os.Create(pathOut + "small-" + f.Name())
		if err != nil {
			log.Fatal(err)
		}

		defer out.Close()
		jpeg.Encode(out, imgOut, nil)

	}

	return
}

func main() {
	//lambda.Start(HandlerRequest)

	req := Request{SourceBucket: "sam", Timeout: 30}
	output, err := HandlerRequest(req)
	if err != nil {
		log.Printf("Error- %v", err)
		return
	}
	log.Print(output)
}
