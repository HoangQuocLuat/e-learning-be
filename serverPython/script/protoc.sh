python -m grpc_tools.protoc \
--python_out=./py/check-in \
--grpc_python_out=./py/check-in \
./proto/check-in.proto