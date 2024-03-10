VERSION=v0.0.1

.PHONY: bin
bin: bin/wplite_darwin_x86_64 bin/wplite_darwin_arm64 bin/wplite_linux_x86_64 bin/wplite_linux_arm bin/wplite_linux_arm64

bin/wplite_darwin_x86_64:
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.Version=$(VERSION)'" -o bin/wplite_darwin_x86_64 cmd/wplite/*.go
	tar -czf bin/wplite_darwin_x86_64.tar.gz bin/wplite_darwin_x86_64
	openssl sha512 bin/wplite_darwin_x86_64 > bin/wplite_darwin_x86_64.sha512
	openssl sha512 bin/wplite_darwin_x86_64.tar.gz > bin/wplite_darwin_x86_64.tar.gz.sha512

bin/wplite_darwin_arm64:
	mkdir -p bin
	GOOS=darwin GOARCH=arm64 go build -ldflags="-X 'main.Version=$(VERSION)'" -o bin/wplite_darwin_arm64 cmd/wplite/*.go
	tar -czf bin/wplite_darwin_arm64.tar.gz bin/wplite_darwin_arm64
	openssl sha512 bin/wplite_darwin_arm64 > bin/wplite_darwin_arm64.sha512
	openssl sha512 bin/wplite_darwin_arm64.tar.gz > bin/wplite_darwin_arm64.tar.gz.sha512

bin/wplite_linux_x86_64:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.Version=$(VERSION)'" -o bin/wplite_linux_x86_64 cmd/wplite/*.go
	tar -czf bin/wplite_linux_x86_64.tar.gz bin/wplite_linux_x86_64
	openssl sha512 bin/wplite_linux_x86_64 > bin/wplite_linux_x86_64.sha512
	openssl sha512 bin/wplite_linux_x86_64.tar.gz > bin/wplite_linux_x86_64.tar.gz.sha512

bin/wplite_linux_arm:
	mkdir -p bin
	GOOS=linux GOARCH=arm go build -ldflags="-X 'main.Version=$(VERSION)'" -o bin/wplite_linux_arm cmd/wplite/*.go
	tar -czf bin/wplite_linux_arm.tar.gz bin/wplite_linux_arm
	openssl sha512 bin/wplite_linux_arm > bin/wplite_linux_arm.sha512
	openssl sha512 bin/wplite_linux_arm.tar.gz > bin/wplite_linux_arm.tar.gz.sha512

bin/wplite_linux_arm64:
	mkdir -p bin
	GOOS=linux GOARCH=arm64 go build -ldflags="-X 'main.Version=$(VERSION)'" -o bin/wplite_linux_arm64 cmd/wplite/*.go
	tar -czf bin/wplite_linux_arm64.tar.gz bin/wplite_linux_arm64
	openssl sha512 bin/wplite_linux_arm64 > bin/wplite_linux_arm64.sha512
	openssl sha512 bin/wplite_linux_arm64.tar.gz > bin/wplite_linux_arm64.tar.gz.sha512