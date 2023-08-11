IMG ?= manager:0.0.1

.PHONY: all
all: help

.PHONY: build
build:
	go build -o bin/manager main.go

.PHONY: run
run:
	go run main.go

.PHONY: docker-build
docker-build:
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push:
	docker push ${IMG}

.PHONY: deploy
deploy:
	kubectl apply -f deploy/manifest.yaml
	kubectl -n mosquitto-operator set image deploy mosquitto-operator *=${IMG}
.PHONY: undeploy
undeploy:
	kubectl delete -f deploy/manifest.yaml --ignore-not-found=true

.PHONY: help
help:
	@echo
	@echo '    * help           Show help'
	@echo '    * build          Build Manager binary'
	@echo '    * run            Run Manager host'
	@echo '    * docker-build   Build Manager image, IMG=manager:0.0.1 (default)'
	@echo '    * docker-push    Push Manager image, IMG=manager:0.0.1 (default)'
	@echo '    * deploy         Deploy Manager in k8s'
	@echo '    * undeploy       Undeploy Manager in k8s'
	@echo
	