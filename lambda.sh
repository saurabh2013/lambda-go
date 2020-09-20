
FUNC_NAME="my-function-sync"

export AWS_PROFILE=testuser
AWS_REGION=us-west-2

source_bucket="testlambdaimages"
target_bucket="testlambdaimages-small"

rebuid(){
echo "=> Re-building code"
rm function.zip main   
GOOS=linux go build main.go resize.go    
zip function.zip main   
}
 
 
if [[ "$1" == "-c" ]] ; then
    echo "=> Creating new lambda func $FUNC_NAME"
    rebuid
    aws lambda create-function --function-name $FUNC_NAME --runtime go1.x \
    --zip-file fileb://function.zip --handler main \
    --role arn:aws:iam::812290705309:role/s3roleforlambda \
    --memory-size 3008 --timeout 900 --region $AWS_REGION \
    --environment Variables={S3_REGION=$AWS_REGION} ;
fi
if [[ "$1" == "-u" ]] ; then
    echo "=> Updating code for $FUNC_NAME"
    rebuid
    aws lambda update-function-code \
    --function-name $FUNC_NAME  --region $AWS_REGION \
    --zip-file fileb://function.zip ;
fi


invoke(){
    echo "=> Function invoke $1"
    aws lambda invoke \
    --function-name $1 --region $AWS_REGION \
    --payload '{"sourcebucket": "'$source_bucket'","destinationbucket": "'$target_bucket'"}' \
    response.json --log-type Tail | jq -r .LogResult| base64 -D;
}

invoke $FUNC_NAME
invoke "my-function-async"

echo "=> Done"
