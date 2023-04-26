image := argocd-autosync-exporter
tag := 0.0.1

build:
	docker build -t $(image):$(tag) -f deploy/docker/dockerfile .
push:
	docker push $(image):$(tag)
release: build push