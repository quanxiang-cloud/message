CONF ?=$(shell pwd)/configs
MESSAGE_SERVER ?=http://127.0.0.1:80
MESSAGE_PORT ?=80
LETTER_PORT ?=8080
EMAIL_PORT ?=8080
PUBSUB_NAME ?=message-kafka-pubsub

REPO ?=dockerhub.qingcloud.com/lowcode
TAG ?=latest

NAMESPACE ?=lowcode

.PHONY: generate
generate:
	go generate ./...

.PHONY: run-letter
run-letter: generate
	dapr run -d ${CONF}/deploy --app-id message-letter -p ${LETTER_PORT} -- go run cmd/letter/main.go --port :${LETTER_PORT} --message-server ${MESSAGE_SERVER}

.PHONY: run-email
run-email: generate
	dapr run -d ${CONF}/deploy --app-id message-letter -p ${EMAIL_PORT} -- go run cmd/email/main.go --port :${EMAIL_PORT} \
		--email-host ${HOST} \
		--email-port ${PORT}\
		--email-username ${USERNAME}\
		--email-password ${PASSWORD}\
		--email-alias ${ALIAS}\
		--email-sender ${SENDER}

.PHONY: run-manager
run-manager: generate
	dapr run -d ${CONF}/deploy --app-id message-manager -p ${MESSAGE_PORT} -- go run cmd/message/main.go --config ${CONF}/config.yml --pubsub-name ${PUBSUB_NAME}

.PHONY: docker-build-letter
docker-build-letter: generate
	KO_DOCKER_REPO=${REPO} ko build -t=${TAG} -B --platform linux/amd64 ./cmd/letter/.

.PHONY: docker-build-email
docker-build-email: generate
	KO_DOCKER_REPO=${REPO} ko build -t=${TAG} -B --platform linux/amd64 ./cmd/email/.

.PHONY: docker-build-message
docker-build-message: generate
	KO_DOCKER_REPO=${REPO} ko build -t=${TAG} -B --platform linux/amd64 ./cmd/message/.

.PHONY: deploy
deploy:
	kubectl apply -n ${NAMESPACE} -f ${CONF}/deploy

.PHONY: undeploy
undeploy:
	kubectl delete -n ${NAMESPACE} -f ${CONF}/deploy