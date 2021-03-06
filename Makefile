# rabtap makefile

BINARY_WIN64=rabtap-win-amd64.exe
BINARY_DARWIN64=rabtap-darwin-amd64
BINARY_LINUX64=rabtap-linux-amd64
SOURCE=$(shell find . -name "*go" -a -not -path "./vendor/*" -not -path "./cmd/testgen/*" )
VERSION=$(shell git describe --tags)
TOXICMD:=docker-compose exec toxiproxy /go/bin/toxiproxy-cli

.PHONY:test-app test-lib build build-all tags short-test test run-server clean \
	   $(BINARY_LINUX64) $(BINARY_WIN64) $(BINARY_DARWIN64)

build:	$(BINARY_LINUX64)

build-all:	build $(BINARY_WIN64)  $(BINARY_DARWIN64)

$(BINARY_DARWIN64): 
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags \
				"-X main.RabtapAppVersion=$(VERSION)" -o $(BINARY_DARWIN64) cmd/main/*.go 

$(BINARY_LINUX64): 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags \
				"-X main.RabtapAppVersion=$(VERSION)" -o $(BINARY_LINUX64) cmd/main/*.go

$(BINARY_WIN64): 
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags \
				"-X main.RabtapAppVersion=$(VERSION)" -o $(BINARY_WIN64) cmd/main/*.go 

tags: $(SOURCE)
	@gotags -f tags $(SOURCE)

lint:
	@./pre-commit

short-test: 
	go test -v -race  github.com/jandelgado/rabtap/cmd/main
	go test -v -race  github.com/jandelgado/rabtap/pkg

test-app:
	go test -race -v -tags "integration" -cover -coverprofile=coverage_app.out github.com/jandelgado/rabtap/cmd/main

test-lib:
	go test -race -v -tags "integration" -cover -coverprofile=coverage.out github.com/jandelgado/rabtap/pkg

test: test-app test-lib
	grep -v "^mode:" coverage_app.out >> coverage.out
	go tool cover -func=coverage.out

toxiproxy-setup:
	$(TOXICMD) c amqp --listen :55672 --upstream rabbitmq:5672 || true

toxiproxy-cmd:
	$(TOXCMD) $(TOXARGS)

# run rabbitmq server for integration test using docker container.
run-broker:
	sudo docker run -ti --rm -p 5672:5672 \
		        -p 15672:15672 rabbitmq:3-management

dist-clean: clean
	rm -f *.out $(BINARY_WIN64) $(BINARY_LINUX64) $(BINARY_DARWIN64)

clean:
	go clean -r

