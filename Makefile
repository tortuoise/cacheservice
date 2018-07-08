#COMMIT = $$(git describe --always)

all:	server client

server:
	@echo "====> Building server"
	go build -o bin/cmd/cacheservice/server github.com/tortuoise/cacheservice/cmd/server

client:
	@echo "====> Building client"
	go build -o bin/cmd/cacheservice/client github.com/tortuoise/cacheservice/cmd/client

deps:
	@echo "====> Install dependencies..."

clean:
	@echo "====> Remove installed binary"
	rm -f bin/aclient

install: deps
	@echo "====> Build aclient in ./bin "
	go build -ldflags "-X main.GitCommit=\"$(COMMIT)\"" -o bin/aclient
