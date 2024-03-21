test:
	go test ./... -coverprofile cover.out.tmp && \
	cat cover.out.tmp | grep -v "IConnector.go" | grep -v "docs.go" > cover.out && \
	rm cover.out.tmp && \
	go tool cover -func cover.out

certs: 
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem