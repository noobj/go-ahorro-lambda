.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/dailyfetch cmd/dailyfetch/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/world cmd/world/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
