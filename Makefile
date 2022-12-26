.PHONY: build clean deploy

swimbuild:
	env GOARCH=amd64 GOOS=linux go build -o bin/dailyfetch cmd/swimming/dailyfetch/main.go
	env GOARCH=amd64 GOOS=linux go build -o bin/crowdscraper cmd/swimming/crowdscraper/main.go
	env GOARCH=amd64 GOOS=linux go build -o bin/swimnotify cmd/swimming/notify/main.go
	env GOARCH=amd64 GOOS=linux go build -o bin/swimnotify_modify cmd/swimming/notification-modify/main.go

build:
	env GOARCH=amd64 GOOS=linux go build -o bin/fetchentries cmd/ahorro/fetchentries/main.go
	env GOARCH=amd64 GOOS=linux go build -o bin/login cmd/ahorro/login/main.go
	env GOARCH=amd64 GOOS=linux go build -o bin/refresh cmd/ahorro/refresh/main.go
	env GOARCH=amd64 GOOS=linux go build -o bin/sync_receiver cmd/ahorro/sync/receiver/main.go
	env GOARCH=amd64 GOOS=linux go build -o bin/sync_callback cmd/ahorro/sync/callback/main.go
	env GOARCH=amd64 GOOS=linux go build -o bin/sync_checker cmd/ahorro/sync/checker/main.go
	env GOARCH=amd64 GOOS=linux go build -o bin/sync_handler cmd/ahorro/sync/handler/*

debug_build:
	env GOARCH=amd64 GOOS=linux go build -gcflags "all=-N -l" -o bin/dailyfetch cmd/swimming/dailyfetch/main.go
	env GOARCH=amd64 GOOS=linux go build -gcflags "all=-N -l" -o bin/crowdscraper cmd/swimming/crowdscraper/main.go
	env GOARCH=amd64 GOOS=linux go build -gcflags "all=-N -l" -o bin/fetchentries cmd/ahorro/fetchentries/main.go
	env GOARCH=amd64 GOOS=linux go build -gcflags "all=-N -l" -o bin/login cmd/ahorro/login/main.go
	env GOARCH=amd64 GOOS=linux go build -gcflags "all=-N -l" -o bin/refresh cmd/ahorro/refresh/main.go
	env GOARCH=amd64 GOOS=linux go build -gcflags "all=-N -l" -o bin/sync_receiver cmd/ahorro/sync/receiver/main.go
	env GOARCH=amd64 GOOS=linux go build -gcflags "all=-N -l" -o bin/sync_callback cmd/ahorro/sync/callback/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
