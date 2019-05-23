build: dep
	@GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"'  -o ./bin/pinger ./cmd/pinger/*.go

dep:
	@GO111MODULE=on go mod vendor -v

run:
	@go run ./cmd/pinger/.

image: build
	@docker build -t zephinzer/go-pinger:latest .