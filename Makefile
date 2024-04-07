run:
	go run ./cmd

test:
	go test -v -coverpkg ./... ./... -coverprofile cover.out.tmp && \
	cat cover.out.tmp | grep -v "mock" | grep -v "docs.go" > cover.out && \
	rm cover.out.tmp && \
	go tool cover -func cover.out

lint:
	golangci-lint run -c linter/.golangci.toml
	
swag:
	swag init -g ./cmd/main.go

MOCKS_DESTINATION=mocks
.PHONY: mocks
mocks: ./internal/service/interfaces.go ./internal/repository/interfaces.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done