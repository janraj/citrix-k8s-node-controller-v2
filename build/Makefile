OWNER=janrajc
IMAGE_NAME=citrix-k8s-node-controller

.PHONY: build
build:
	docker build -f Dockerfile -t $(IMAGE_NAME) .. -m 4g
	docker save -o $(IMAGE_NAME).tar $(IMAGE_NAME)
test:
	 go test -v ../... -coverprofile cp.out 
	 go tool cover -html=cp.out -o cover.html

coverage:
	 go test -v ../... -coverprofile=coverage.txt -covermode=atomic

publish: repo-login publish

repo-login:
	@docker login -u "$(QUAY_USERNAME)" -p "$(QUAY_PASSWORD)" quay.io

publish:
	@echo 'publish latest to $(DOCKER_REGISTRY)'
	docker tag  $(IMAGE_NAME):latest $(DOCKER_REGISTRY)/$(IMAGE_NAME):latest
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):latest

clean:
	docker rmi -f $$(docker images -q -f dangling=true) || true
	docker rmi -f $$(docker images | awk '$$1 ~ /citrix-k8s-node-controller/ { print $$3}') || true
         
