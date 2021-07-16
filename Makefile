SHELL := bash

image:
	mkdir -p build/bin
	find testcase/ -name '*.test' | xargs -I{} cp {} build/bin/
	docker build -t hub.pingcap.net/cosven/e2etest build/ -f Dockerfile

.PHONY: build
