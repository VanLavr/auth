run:
	@docker-compose up --build -d

logs:
	@docker logs auth-authservice-1

stoprm:
	@docker stop auth-authservice-1 auth-tokens_db-1
	@docker rm auth-authservice-1 auth-tokens_db-1

# mockery \
# > --dir=internal/auth/service \
# > --filename=auth_repo_mocks.go \
# > --name=Repository \
# > --output=internal/mocks/auht/repo \
# > --outpkg=auth_repo_mocks