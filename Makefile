.PHONY: terraform
terraform:
	terraform -chdir=terraform/local/ init
	terraform -chdir=terraform/local/ plan
	terraform -chdir=terraform/local/ apply -auto-approve

.PHONY: drop-database
drop-database:
	psql -h localhost -U postgres -d postgres -c "DROP DATABASE IF EXISTS mucaron WITH (FORCE);"
	psql -h localhost -U postgres -d postgres -c "CREATE DATABASE mucaron;"

# .PHONY: mkcert
# mkcert:
# 	mkcert -install
# 	mkcert -cert-file ./certs/mucaron.local.walnuts.dev.pem -key-file ./certs/mucaron.local.walnuts.dev-key.pem mucaron.local.walnuts.dev
# 	mkcert -cert-file ./certs/minio.local.walnuts.dev.pem -key-file ./certs/minio.local.walnuts.dev-key.pem minio.local.walnuts.dev

certs/mucaron.local.walnuts.dev.pem:
	mkcert -cert-file ./certs/mucaron.local.walnuts.dev.pem -key-file ./certs/mucaron.local.walnuts.dev-key.pem mucaron.local.walnuts.dev
	mkcert -install

certs/minio.local.walnuts.dev.pem:
	mkcert -cert-file ./certs/minio.local.walnuts.dev.pem -key-file ./certs/minio.local.walnuts.dev-key.pem minio.local.walnuts.dev
	mkcert -install

.PHONY: up
up: certs/mucaron.local.walnuts.dev.pem certs/minio.local.walnuts.dev.pem
	docker-compose up -d
	cd frontend && yarn install && yarn dev


.PHONY: down
down:
	docker-compose down
