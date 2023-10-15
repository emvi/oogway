.PHONY: deps test build_mac build_windows build_linux

deps:
	go get -u -t ./...
	go mod vendor

test:
	go test -cover ./pkg

build_mac: test
	GOOS=darwin go build -a -installsuffix cgo -ldflags "-s -w" cmd/oogway/main.go

build_windows: test
	GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" -o oogway.exe cmd/oogway/main.go

build_linux: test
	GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -w" cmd/oogway/main.go
