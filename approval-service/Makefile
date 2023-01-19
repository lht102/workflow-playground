SHELL := /bin/sh

MIGRATION_PATH := ./migrations/

-include .env.$(ENV)
export

.PHONY: run-api-server
run-api-server:
	go run cmd/api-server/main.go

.PHONY: run-worker
run-worker:
	go run cmd/worker/main.go

.PHONY: dep
dep:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run

.PHONY: build
build:
	@rm -rf bin/
	@mkdir bin/
	@cd bin && go build -v -o . ../cmd/api-server
	@cd bin && go build -v -o . ../cmd/worker

.PHONY: test
test:
	go test ./... -race -cover -coverprofile=c.out -covermode=atomic

.PHONY: mocks
mocks:
	mockery --name=PaymentService --with-expecter

.PHONY: gen-api-server
gen-api-server:
	oapi-codegen -package api -generate types,chi-server,spec api/openapi.yaml > api/api.gen.go

.PHONY: create-model
create-model:
	GOWORK=off go run -mod=mod entgo.io/ent/cmd/ent init $(name)

.PHONY: gen
gen:
	GOWORK=off go generate ./...

.PHONY: migrate-create
migrate-create:
	go run cmd/migrate/main.go $(name)

.PHONY: migrate-manual-create
migrate-manual-create:
	@mkdir -p $(MIGRATION_PATH)
	@migrate create -dir $(MIGRATION_PATH) -ext sql $(name)
	@atlas migrate hash --dir file://migrations

.PHONY: migrate-validate
migrate-validate:
	atlas migrate validate

.PHONY: migrate
migrate:
	@migrate -path $(MIGRATION_PATH) \
		-database "mysql://$(MYSQL_USERNAME):$(MYSQL_PASSWORD)@$(MYSQL_PROTOCOL)($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" up

.PHONY: migrate-rollback
migrate-rollback:
	@migrate -path $(MIGRATION_PATH) \
		-database "mysql://$(MYSQL_USERNAME):$(MYSQL_PASSWORD)@$(MYSQL_PROTOCOL)($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" down 1

.PHONY: migrate-drop
migrate-drop:
	@migrate -path $(MIGRATION_PATH) \
		-database "mysql://$(MYSQL_USERNAME):$(MYSQL_PASSWORD)@$(MYSQL_PROTOCOL)($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" down

.PHONY: reset-db
reset-db:
	@mysql -h $(MYSQL_ADDRESS) -u $(MYSQL_USERNAME) --password=$(MYSQL_PASSWORD) --protocol=$(MYSQL_PROTOCOL) -e "DROP DATABASE IF EXISTS $(MYSQL_DATABASE); CREATE DATABASE $(MYSQL_DATABASE);"

.PHONY: migrate-reset
migrate-reset: reset-db migrate
