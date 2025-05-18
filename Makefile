SHELL := /bin/bash

curr_dir := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
rest_args := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))
$(eval $(rest_args):;@:)

help:
	#
	# Usage:
	#   make generate :  generate a mapper based on a template.
	#   make mapper {mapper-name} <action> <parameter>:  execute mapper building process.
	#
	# Actions:
	#   -           mod, m  :  download code dependencies.
	#   -          lint, l  :  verify code via go fmt and `golangci-lint`.
	#   -         build, b  :  compile code.
	#   -       package, p  :  package docker image.
	#   -         clean, c  :  clean output binary.
	#
	# Parameters:
	#   ARM   : true or undefined
	#   ARM64 : true or undefined
	#
	# Example:
	#   -  make mapper modbus ARM64=true :  execute `build` "modbus" mapper for ARM64.
	#   -        make mapper modbus test :  execute `test` "modbus" mapper.
	@echo

make_rules := $(shell ls $(curr_dir)/hack/make-rules | sed 's/.sh//g')
$(make_rules):
	@$(curr_dir)/hack/make-rules/$@.sh $(rest_args)

.DEFAULT_GOAL := help
.PHONY: $(make_rules) build test package


DOCKER_REGISTRY ?= ""

deploy: deploy-crds deploy-resource

clean:  undeploy-crds undeploy-resource
build-app:
	CGO_ENABLED=0 GOOS=linux go build  -ldflags="-s -w" -o main ./cmd/main.go
docker-build-push:
ifeq ($(DOCKER_REGISTRY), "")
	$(error DOCKER_REGISTRY is not set, please set it use "export DOCKER_REGISTRY=<your-registry> " first)
endif
	docker buildx build -f ./Dockerfile_nostream -t ${DOCKER_REGISTRY}/temperature-mapper:v1.0 .
	docker push ${DOCKER_REGISTRY}/temperature-mapper:v1.0
undeploy-crds:
	kubectl delete -f./crds/temperature-instance.yaml
	kubectl delete -f./crds/temperature-model.yaml
undeploy-resource:
	kubectl delete -f./resource/deployment.yaml
	kubectl delete -f./resource/configmap.yaml

deploy-resource:
	kubectl apply -f ./resource/configmap.yaml
	kubectl apply -f./resource/deployment.yaml
deploy-crds:
	kubectl apply -f./crds/temperature-model.yaml
	kubectl apply -f./crds/temperature-instance.yaml

