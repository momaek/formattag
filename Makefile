Version=$$(git describe --tags)
build_macOS_x86:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=${Version}" -o formattag-macOS-amd64 main.go

build_macOS_arm:
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=${Version}" -o formattag-macOS-arm64 main.go

build_macOS: build_macOS_arm build_macOS_x86

build_linux_x86:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=${Version}" -o formattag-linux-amd64 main.go

build_linux_arm64:
	GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=${Version}"  -o formattag-linux-arm64 main.go

build_linux_386:
	GOOS=linux GOARCH=386 go build -ldflags "-X main.version=${Version}" -o formattag-linux-386 main.go

build_linux: build_linux_386 build_linux_arm64 build_linux_x86 

build: build_linux build_macOS

tar:
	tar zcvf formattag.tar.gz formattag-macOS-amd64 formattag-macOS-arm64  formattag-linux-amd64 formattag-linux-arm64 formattag-linux-386
