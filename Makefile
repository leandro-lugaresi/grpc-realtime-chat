.PHONY: generate-protobuf

#### PROJECT SETTINGS ####
current_path = $(dir $(mkfile_path))
server_path = server
export_go_path = export GOPATH=$(current_path)/vendor:$(current_path)

#### END PROJECT SETTINGS ####

#### TARGETS ####

all: generate-protobuf

download-protofuf:
	rm -rf /tmp/googleapis && \
	git clone https://github.com/google/googleapis.git /tmp/googleapis

generate-protobuf:
	protoc -I./proto -I /tmp/googleapis ./proto/chat.proto --go_out=plugins=grpc:$(server_path)/chat && \
	protoc -I./proto -I /tmp/googleapis ./proto/user.proto --go_out=plugins=grpc:$(server_path)/user
