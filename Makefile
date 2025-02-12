
rpc-service:
	go build -v $(LDFLAGS)

clean:
	rm rpc-service

test:
	go test -v ./...

lint:
	golangci-lint run ./...

.PHONY: \
	rpc-service \
	clean \
	test \
	lint