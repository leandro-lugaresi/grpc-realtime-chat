# Instalation

- `npm install`

## Update client protobuf
Code for this client is pre-generated using protoc and the Node gRPC protoc plugin, and the generated code can be found in various `*_pb.js` files.
The command line sequence for generating those files is as follows (assuming that `protoc` and `grpc_node_plugin` are present, and starting in the base directory of this package):

```sh
cd ../protos
node_modules/grpc_tools_node_protoc --js_out=import_style=commonjs,binary:. --grpc_out=. --plugin=protoc-gen-grpc=node_modules/.bin/grpc_tools_node_protoc_plugin
protoc --js_out=import_style=commonjs,binary:../node/static_codegen/route_guide/ --grpc_out=../node/static_codegen/route_guide/ --plugin=protoc-gen-grpc=grpc_node_plugin route_guide.proto
```

