VERSION := 0.0.4

kind-up: 
	kind create cluster \
		--image kindest/node:v1.29.0 \
		--name phalcon-test

deploy-nats:
	helm install nats nats/nats

ping:
	docker build \
		-f zarf/docker/dockerfile.ping \
		-t jnkroeker/ping-image:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

pong:
	docker build \
		-f zarf/docker/dockerfile.pong \
		-t jnkroeker/pong-image:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

outliers:
	docker build \
		-f zarf/docker/dockerfile.outliers \
		-t jnkroeker/outliers-image:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

kind-load:
	kind load docker-image jnkroeker/ping-image:$(VERSION) --name phalcon-test
	kind load docker-image jnkroeker/pong-image:$(VERSION) --name phalcon-test

kind-ping-apply:
	kubectl apply -f zarf/k8s/base-ping.yaml 

kind-pong-apply:
	kubectl apply -f zarf/k8s/base-pong.yaml

ping-logs:
	kubectl logs -l app=ping --namespace=default --all-containers=true -f --tail=100

pong-logs:
	kubectl logs -l app=pong --namespace=default --all-containers=true -f --tail=100

proto:
	python3 -m grpc_tools.protoc -I ./protobufs --python_out=. --grpc_python_out=. ./protobufs/outliers.proto