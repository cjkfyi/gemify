#!/bin/bash

echo -e "\n🚀  Launching 'scripts/gen.sh'\n" 

MODULE="gemify/api"
GO_OUT_DIR="backend/api/gen"
PROTO_FILE="backend/api/proto/api.proto"

rm -rf ${GO_OUT_DIR}/* 
echo -e "💦  Cleaned up prev generations!\n" 

echo -e "🖨️   Starting 'protoc' codegen!\n" 
protoc \
--go_out=${GO_OUT_DIR} \
--go_opt=module=${MODULE} \
--go-grpc_out=${GO_OUT_DIR} \
--go-grpc_opt=module=${MODULE} \
${PROTO_FILE}
echo -e "✨  Generated our Go files!\n" 

echo -e "👻  Ending 'scripts/gen.sh'\n" 
