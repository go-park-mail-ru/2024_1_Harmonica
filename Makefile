test:
	go test -coverpkg=./... -coverprofile=cover.out ./... && go tool cover -func cover.out | grep total: