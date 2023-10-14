.PHONY: gen-fibonacci-proto 
gen-fibonacci-proto: 
	protoc --go_out=. --go-grpc_out=. proto/*.proto