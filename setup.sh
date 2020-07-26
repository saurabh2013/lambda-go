#!/bin/bash
echo "GOPATH: " $(go env GOPATH)
echo "GOROOT: " $(go env GOROOT)
 
GO111MODULE=on go get -v github.com/urfave/cli/v2
go get -v github.com/stretchr/testify/assert
go get -v golang.org/x/net/http2
go get -v github.com/aws/aws-lambda-go/lambda
go get -v github.com/aws/aws-sdk-go




