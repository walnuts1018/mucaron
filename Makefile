.PHONY: terraform
terraform:
	terraform -chdir=terraform/local/ init
	terraform -chdir=terraform/local/ plan
	terraform -chdir=terraform/local/ apply -auto-approve

