if [[ "$PWD" =~ src ]]; then
export GOPATH="$(dirname $(dirname "$(pwd)"))"
else 
export GOPATH=$(pwd)
fi
echo "Updated GOPATH: " $(go env GOPATH)
export PATH=$GOPATH/bin:$PATH


source setup.sh
go run main.go resize.go
