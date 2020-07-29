
FUNC_NAME="my-function-async1"
export AWS_PROFILE=test

source_bucket="testlambdaimages"
target_bucket="testlambdaimages-small"

rebuid(){
echo "=> Re-building code"
rm function.zip main   
GOOS=linux go build main.go resize.go    
zip function.zip main   
}
 
 
if [[ "$1" == "c" ]] ; then
    echo "=> Creating new lambda func $FUNC_NAME"
    rebuid
    aws lambda create-function --function-name $FUNC_NAME --runtime go1.x \
    --zip-file fileb://function.zip --handler main \
    --role arn:aws:iam::812290705309:role/s3roleforlambda \
    --memory-size 3008 --timeout 900 \
    --environment Variables={S3_REGION=us-west-2} ;
fi
if [[ "$1" == "u" ]] ; then
    echo "=> Updating code for $FUNC_NAME"
    rebuid
    aws lambda update-function-code \
    --function-name $FUNC_NAME \
    --zip-file fileb://function.zip ;
fi


invoke(){
    echo "=> Function invoke $1"
    aws lambda invoke \
    --function-name $1 \
    --payload '{"sourcebucket": "'$source_bucket'","destinationbucket": "'$target_bucket'"}' \
    response.json --log-type Tail | jq -r .LogResult| base64 -D;
}

invoke $FUNC_NAME
invoke "my-function-sync"

echo "=> Done"
