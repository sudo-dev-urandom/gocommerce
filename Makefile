BIN_NAME=gocommerce
#!make
include .env
export $(shell sed 's/=.*//' .env)

exports:
	@printenv | grep MYAPP

run: exports generate-docs
	@go run main.go

generate-docs:
	@echo "Updating API documentation..."
	@swag init


unit-test:
	@go test tests/unit/*_test.go