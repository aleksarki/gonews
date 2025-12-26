@echo off

set PROTO_DIR=protos
set OUTPUT_DIR=protos\pb

if not exist %OUTPUT_DIR% mkdir %OUTPUT_DIR%

protoc --go_out=%OUTPUT_DIR% ^
       --go_opt=paths=source_relative ^
       --go-grpc_out=%OUTPUT_DIR% ^
       --go-grpc_opt=paths=source_relative ^
       -I %PROTO_DIR% ^
       %PROTO_DIR%\*.proto
