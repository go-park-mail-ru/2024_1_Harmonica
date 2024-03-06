test:
	go test ./... -coverprofile cover.out.tmp && \
	cat cover.out.tmp | grep -v "Methods.go" > cover.out && \
	rm cover.out.tmp && \
	go tool cover -func cover.out