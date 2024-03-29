rm -rf ./proto_codes
mkdir ./proto_codes
mkdir ./proto_codes/db
mkdir ./proto_codes/rpc
protoc --go_out=. ./proto/db_proto/*.proto
protoc --go_out=. ./proto/rpc_proto/*.proto