.PHONY:build
build:
	# build app for linux based container
	@CGO_ENABLED=0 GOOS=linux go build -o ./fizzbuzz -a -ldflags '-s' -installsuffix cgo ./cmd/fizzbuzz-api/main.go

.PHONY:test
test:
	# lauch unit test across all project
	@go test -race ./...
	@go list ./... | grep -v /vendor/ | grep -v pb | xargs -L1 golint -set_exit_status
	@go vet `go list ./... | grep -v /vendor/`

.PHONY:coverage
coverage:
	# use go ability to generate an html page with test coverage
	@go test `go list ./... | grep -v /vendor/` -cover -coverprofile=cover.out
	@go tool cover -html=cover.out

.PHONY:dev
dev:
	# launch external dependencies from docker-compose for local development
	@docker-compose up db jaeger

.PHONY:local
local:
	# launch complete docker-compose for local execution
	@docker-compose up

.PHONY:clean
clean:
	# clean docker compose
	@docker-compose down
	@docker rmi fizzbuzz

.PHONY:integration
integration:
	# launch integration test
	@CGO_ENABLED=0 go test -v `go list ./... | grep -v /vendor/` -tags=integration

.PHONY: run
run:
	@go run cmd/fizzbuzz-api/main.go
