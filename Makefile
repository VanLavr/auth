run:
	@go run cmd/main.go

gen_mocks:
	@mockery --dir=./internal/auth/service --name=Repository