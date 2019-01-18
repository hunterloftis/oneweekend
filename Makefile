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
	
profile:
	go build ./cmd/trace
	./trace -profile > test.ppm
	@echo 'Next, go tool pprof --pdf ./trace /tmp/path/to/cpu.pprof > cpu.pdf'

profwin:
	PATH %PATH%;C:\Program Files (x86)\Graphviz2.38\bin
	go build ./cmd/trace
	trace.exe -profile > test.ppm
	go tool pprof --pdf ./trace.exe C:\path\to\cpu.pprof > cpu.pdf
	start cpu.pdf
