.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build binary
	GOOS=linux GOARCH=amd64 go build -o elevenbot main.go

deploy: ## Deploys binary on a Google Cloud instance
	@echo '>>> Deploy started'

	ssh ${ELEVENBOT_SSH_USER}@${ELEVENBOT_SSH_IP} 'sudo systemctl stop elevenbot'
	scp elevenbot ${ELEVENBOT_SSH_USER}@${ELEVENBOT_SSH_IP}:~
	ssh ${ELEVENBOT_SSH_USER}@${ELEVENBOT_SSH_IP} 'sudo systemctl start elevenbot'

	@echo '>>> done!'

prepare-deploy: ## Prepares bot to be deployed on a new instance
	@echo '>>> Uploads config/ to server home directory'

	scp -r config/ ${ELEVENBOT_SSH_USER}@${ELEVENBOT_SSH_IP}:~
	ssh ${ELEVENBOT_SSH_USER}@${ELEVENBOT_SSH_IP} 'sudo mv config/elevenbot.service /etc/systemd/system/elevenbot.service'
	ssh ${ELEVENBOT_SSH_USER}@${ELEVENBOT_SSH_IP} 'sudo systemctl daemon-reload'

	@echo '>>> done!'

install: ## Installs dependencies
	@echo '>>> Installs golang dependencies'

	dep ensure

	@echo '>>> done!'
