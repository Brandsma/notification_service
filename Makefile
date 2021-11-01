.PHONY: FORCE

IMAGE_NAME=notification
REGISTRY=europe-west4-docker.pkg.dev/wacc-schrobox/schrobox-registry/${IMAGE_NAME}

run-dev: FORCE guard-REGISTRY
	echo "running: ${REGISTRY}:dev"
	docker run --env-file ../.env -p 9003:9003 -t ${REGISTRY}:dev .

build-dev: FORCE guard-REGISTRY
	echo "building: ${REGISTRY}:dev"
	docker build --ssh default -t ${REGISTRY}:dev .

build-release: FORCE guard-VERSION guard-REGISTRY
	echo "building: ${REGISTRY}:${VERSION}"
	docker build --ssh default -t ${REGISTRY}:${VERSION} .
	docker push ${REGISTRY}:${VERSION}

guard-%:
	@ if [ -z '${${*}}' ]; then echo 'Environment variable $* not set' && exit 1; fi
