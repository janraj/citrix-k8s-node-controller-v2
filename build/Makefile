OWNER=Janrajc
IMAGE_NAME=citrix-k8s-node-controller
VERSION_FILE="../version/VERSION"
version=0.0.0
error=0.0.0


.PHONY: build
build:
	docker build -f Dockerfile -t $(IMAGE_NAME) .. -m 4g
	docker save -o $(IMAGE_NAME).tar $(IMAGE_NAME)
test:
	@echo 'Commit Message'
	@echo $(TRAVIS_COMMIT_MESSAGE)
	go test -v ../... -coverprofile cp.out 
	go tool cover -html=cp.out -o cover.html

run:
	cd ../cmd/citrix-node-controller/;go run .
coverage:
	go test -v ../... -coverprofile=coverage.txt -covermode=atomic

publish: repo-login version-update release

repo-login:
	@docker login -u "$(QUAY_USERNAME)" -p "$(QUAY_PASSWORD)" quay.io

version-update:
	@echo 'Updating Version'
	$(eval version = $(shell ./gitpush.sh | tail -1))
	@echo 'New versions $(version)'
release:
	@echo 'publish latest and $(version) to $(DOCKER_REGISTRY)'
ifneq ($(version), $(error))
	docker tag  $(IMAGE_NAME):latest $(DOCKER_REGISTRY)/$(IMAGE_NAME):latest
	docker tag  $(IMAGE_NAME):latest $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(version)
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):latest
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(version)
endif

clean:
	docker rmi -f $$(docker images -q -f dangling=true) || true
	docker rmi -f $$(docker images | awk '$$1 ~ /citrix-k8s-node-controller/ { print $$3}') || true
         
