CONF ?=$(shell pwd)/configs
MESSAGE_SERVER ?=http://127.0.0.1:80
MESSAGE_PORT ?=80
LETTER_PORT ?=8080
EMAIL_PORT ?=8080
PUBSUB_NAME ?=message-kafka-pubsub

REPO ?=dockerhub.qingcloud.com/lowcode
TAG ?=latest

NAMESPACE ?=lowcode

generate:
	go generate ./...
run-letter: generate
	dapr run -d ${CONF}/deploy --app-id message-letter -p ${LETTER_PORT} -- go run cmd/letter/main.go --port :${LETTER_PORT} --message-server ${MESSAGE_SERVER}

run-email: generate
	dapr run -d ${CONF}/deploy --app-id message-letter -p ${EMAIL_PORT} -- go run cmd/email/main.go --port :${EMAIL_PORT} \
		--email-host ${HOST} \
		--email-port ${PORT}\
		--email-username ${USERNAME}\
		--email-password ${PASSWORD}\
		--email-alias ${ALIAS}\
		--email-sender ${SENDER}

run-manager: generate
	dapr run -d ${CONF}/deploy --app-id message-manager -p ${MESSAGE_PORT} -- go run cmd/message/main.go --config ${CONF}/config.yml --pubsub-name ${PUBSUB_NAME}

docker-build-letter: generate
	KO_DOCKER_REPO=${REPO} ko build -t=${TAG} -B --platform linux/amd64 ./cmd/letter/.

docker-build-email: generate
	KO_DOCKER_REPO=${REPO} ko build -t=${TAG} -B --platform linux/amd64 ./cmd/email/.

docker-build-message: generate
	KO_DOCKER_REPO=${REPO} ko build -t=${TAG} -B --platform linux/amd64 ./cmd/message/.

deploy:
	kubectl apply -n ${NAMESPACE} -f ${CONF}/deploy

undeploy:
	kubectl delete -n ${NAMESPACE} -f ${CONF}/deploy