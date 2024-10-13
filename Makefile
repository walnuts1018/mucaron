.PHONY: terraform
terraform:
	terraform -chdir=terraform/local/ init
	terraform -chdir=terraform/local/ plan
	terraform -chdir=terraform/local/ apply -auto-approve

.PHONY: drop-database
drop-database:
	psql -h localhost -U postgres -d postgres -c "DROP DATABASE IF EXISTS mucaron WITH (FORCE);"
	psql -h localhost -U postgres -d postgres -c "CREATE DATABASE mucaron;"
