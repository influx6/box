Cmd=./cmd/box
Example=./example
Linuxenv= GOOS=linux GOARCH=amd64 CGO_ENABLED=0
Archenv= GOOS=linux GOARCH=amd64 CGO_ENABLED=0
Osxenv= GOOS=darwin GOARCH=x86_64 CGO_ENABLED=0

build:
	go get -v
	$(Linuxenv) go build -o ./example/ubuntu/box $(Cmd)

install: build
	go install ./cmd/...

docker: build install
	docker build -t ubuntu-run $(Example)/ubuntu/.

clean:
	rm $(Example)/ubuntu/box
