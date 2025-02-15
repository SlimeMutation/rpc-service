GITCOMMIT := $(shell git rev-parse HEAD)
GITDATE := $(shell git show -s --format='%ct')

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

go-signature:
	go build -v $(LDFLAGS) ./cmd/go-signature

clean:
	rm go-signature

test:
	go test -v ./...

lint:
	golangci-lint run ./...

.PHONY: \
	go-signature \
	clean \
	test \
	lint