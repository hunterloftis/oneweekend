test:
	go run cmd/trace/main.go > test.ppm
	open test.ppm

doc:
	mkdir -p /tmp/goroot/doc
	rm -rf /tmp/gopath/src/github.com/hunterloftis/oneweekend
	mkdir -p /tmp/gopath/src/github.com/hunterloftis/oneweekend
	tar -c --exclude='.git' --exclude='tmp' . | tar -x -C /tmp/gopath/src/github.com/hunterloftis/oneweekend
	echo -e "open http://localhost:6060/pkg/github.com/hunterloftis/oneweekend\n"
	GOROOT=/tmp/goroot/ GOPATH=/tmp/gopath/ godoc -http=localhost:6060
	