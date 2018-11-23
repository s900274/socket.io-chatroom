export GO111MODULE=on

all: clean bld

bld: build-service

build-service:
	go build -o bin/socket-server ./src/server

clean:
	@rm -f bin/socket-server

cleanlog:
	@rm -f log/*log*