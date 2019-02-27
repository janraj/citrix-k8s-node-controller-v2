.PHONY: build
build:
	docker build -f Dockerfile -t citrix-node-controller .. -m 4g
	docker save -o citrix-node-controller.tar citrix-node-controller
test:
	 go test ../... -coverprofile cp.out 
	 go tool cover -html=cp.out -o cover.html
clean:
	docker rmi -f $$(docker images -q -f dangling=true) || true
	docker rmi -f $$(docker images | awk '$$1 ~ /citrix-node-controller/ { print $$3}') || true
         
