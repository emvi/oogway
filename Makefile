.PHONY: deps test build_mac build_windows build_linux

deps:
	go get -u -t ./...
	go mod vendor

test:
	go test -cover ./oogway

build_mac: test
	GOOS=darwin go build -a -installsuffix cgo -ldflags "-s -w" -o oogway cmd/main.go

build_windows: test
	GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -o oogway.exe cmd/main.go

build_linux: test
	GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" cmd/main.go
