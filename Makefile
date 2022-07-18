.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/dailyfetch cmd/swimming/dailyfetch/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/crowdscraper cmd/swimming/crowdscraper/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/fetchentries cmd/ahorro/fetchentries/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
