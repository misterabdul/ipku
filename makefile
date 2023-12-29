PORT ?= 3000
BEHIND_PROXY ?= false

.PHONY: all run

all: bin/ipku	

run: bin/ipku
	./bin/ipku -port=$(PORT) -behind-proxy=$(BEHIND_PROXY)

bin/ipku:
	CGO_ENABLED=0 go build -o bin/ipku src/*.go
