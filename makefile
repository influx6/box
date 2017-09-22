CMD=./cmd/box
EXAMPLES=./example
LINUXENV= GOOS=linux GOARCH=amd64 CGO_ENABLED=0
ARCHENV= GOOS=linux GOARCH=amd64 CGO_ENABLED=0
OSXHENV= GOOS=darwin GOARCH=x86_64 CGO_ENABLED=0

build:
	go get -v
	$(LINUXENV) go build -o ./example/ubuntu/box $(CMD)
	$(install)

install:
	go install ./cmd/...

clean:
	rm $(EXAMPLES)/ubuntu/box


$(build)
