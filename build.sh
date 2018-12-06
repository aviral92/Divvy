cd src/pb/
protoc --go_out=plugins=grpc:. divvy.proto
cd -
cd src/node/
go get -v ./...
go build .
cd -
