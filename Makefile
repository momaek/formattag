build_x86:
	GOOS=darwin GOARCH=amd64 go build -o formattag-amd64 main.go

build_m1:
	GOOS=darwin GOARCH=arm64 go build -o formattag-arm64 main.go

build: build_x86 build_m1

tar:
	tar zcvf formattag.tar.gz formattag-amd64 formattag-arm64 
