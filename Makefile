PREFIX=/usr/local

all: build

deps:
	./build.sh --deps

build: deps
	./build.sh

#install: build
# 	/bin/mv mqtt_stresser $(PREFIX)/bin

test: deps
	./run_all_tests.sh
