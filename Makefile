run:
	@docker-compose up --build -d

logs:
	@docker logs auth-authservice-1

stoprm:
	@docker stop auth-authservice-1 auth-tokens_db-1
	@docker rm auth-authservice-1 auth-tokens_db-1