export PATH="$PATH:$(go env GOPATH)/bin"
ls -l $(go env GOPATH)/bin/protoc-gen*

// protoc   --proto_path=.   --go_out=.   --go-grpc_out=.   chord.proto