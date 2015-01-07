all: clean deps build

build:
	go build -v -o $(GOPATH)/bin/packer-provisioner-wait provisioner.go

test:
	go test

clean:; rm -f $(GOPATH)/bin/packer-provisioner-wait

deps:
	go get -d -v -p 2 ./...
