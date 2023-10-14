.PHONY: gen-proto 
gen-proto: 
	protoc --go_out=. --go-grpc_out=. proto/*.proto