.PHONY: build
build:
	docker build -f Dockerfile -t citrix-k8s-node-controller .. -m 4g
	docker save -o citrix-k8s-node-controller.tar citrix-k8s-node-controller
test:
	 go test ../... -coverprofile cp.out 
	 go tool cover -html=cp.out -o cover.html

publish: repo-login publish

repo-login:
	@docker login -u "$(QUAY_USERNAME)" -p "$(QUAY_PASSWORD)" quay.io

publish-latest:
	@echo 'publish latest to $(DOCKER_REGISTRY)'
	docker tag citrix-k8s-node-controller:latest $(DOCKER_REGISTRY)/citrix-k8s-node-controller:latest
	docker push $(DOCKER_REGISTRY)/citrix-k8s-node-controller:latest

clean:
	docker rmi -f $$(docker images -q -f dangling=true) || true
	docker rmi -f $$(docker images | awk '$$1 ~ /citrix-k8s-node-controller/ { print $$3}') || true
         
