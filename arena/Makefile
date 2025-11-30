GO ?= go
GOEXPERIMENT ?= arenas
APP := ./cmd/arena-demo

.PHONY: run fmt tidy

run:
	GOEXPERIMENT=$(GOEXPERIMENT) $(GO) run $(APP)

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy
