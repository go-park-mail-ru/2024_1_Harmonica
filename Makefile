run:
	go run ./cmd/image/main.go & \
	go run ./cmd/like/main.go & \
	go run ./cmd/auth/main.go & \
	go run ./cmd/main.go &

stop:
	for port in ':8002', ':8003', ':8004', ':8080' ; do \
    	for pid in `lsof -i $$port | awk '{print $$2}'`; do \
    		if [ "$$pid" = "PID" ] ; then \
    			continue ; \
    		fi ; \
    		echo $$pid ; \
    		`kill $$pid` ; \
    	done ; \
    done

deploy:
	make stop ; \
	git stash; \
	git checkout main; \
	git pull; \
	nohup make run 

test:
	go test -v -coverpkg ./... ./... -coverprofile cover.out.tmp && \
	cat cover.out.tmp | grep -v "mock" | grep -v "docs.go" | grep -v "proto" > cover.out && \
	rm cover.out.tmp && \
	go tool cover -func cover.out

lint:
	golangci-lint run -c linter/.golangci.toml --skip-dirs='(tests)'
	
swag:
	swag init -g ./cmd/main.go

MOCKS_DESTINATION=mocks
.PHONY: mocks
mocks:  ./internal/service/interfaces.go \
		./internal/repository/interfaces.go \
		./internal/microservices/image/proto/image_grpc.pb.go \
		./internal/microservices/auth/proto/auth_grpc.pb.go \
		./internal/microservices/like/proto/like_grpc.pb.go \
		./internal/microservices/like/server/service/interfaces.go \
		./internal/microservices/like/server/repository/interfaces.go \
		./internal/microservices/auth/server/service/interfaces.go \
		./internal/microservices/auth/server/repository/interfaces.go \
		./internal/microservices/image/server/service/interfaces.go \
		./internal/microservices/image/server/repository/interfaces.go 
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do \
        		dest=$(MOCKS_DESTINATION)/$$(echo $$file | cut -c 10-); \
        		mkdir -p $$(dirname $$dest); \
        		mockgen -source=$$file -destination=$$dest; done

proto_auth:
	export PATH="$(PATH):$(go env GOPATH)/bin" && \
	protoc --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative ./internal/microservices/auth/proto/auth.proto

proto_image:
	export PATH="$(PATH):$(go env GOPATH)/bin" && \
	protoc --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative ./internal/microservices/image/proto/image.proto

proto_like:
	export PATH="$(PATH):$(go env GOPATH)/bin" && \
	protoc --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative ./internal/microservices/like/proto/like.proto

run_auth:
	go run cmd/auth/main.go

run_image:
	go run cmd/image/main.go

run_likes:
	go run cmd/like/main.go
