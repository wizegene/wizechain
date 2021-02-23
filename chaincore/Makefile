generate:

		protoc --proto_path=${GOPATH}/src \
        --proto_path=`go list -f '{{ .Dir }}' -m github.com/mwitkow/go-proto-validators` \
        --proto_path=`go list -f '{{ .Dir }}' -m github.com/grpc-gateway/v2/protoc-gen-openapiv2` \
        --proto_path=./proto \
        --gogo_out=./core/wire \
        --go-grpc_out=./core/wire \
        --govalidators_out=./core/wire \
        ./proto/*.proto

		statik -m -f -src third_party/OpenAPI/



install:
	go get \
		github.com/gogo/protobuf/protoc-gen-gogo \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
		github.com/mwitkow/go-proto-validators/protoc-gen-govalidators \
		github.com/rakyll/statik