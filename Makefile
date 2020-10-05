WEBHOOK_SERVICE?=hello-webhook-service
NAMESPACE?=default
CONTAINER_REPO?=quay.io/didil/hello-webhook
CONTAINER_VERSION?=0.1.9
CONTAINER_IMAGE=$(CONTAINER_REPO):$(CONTAINER_VERSION)

docker-build:
	docker build -t $(CONTAINER_IMAGE) webhook

docker-push:
	docker push $(CONTAINER_IMAGE) 


.PHONY: k8s-deploy
k8s-deploy: k8s-deploy-other k8s-deploy-csr k8s-deploy-deployment

.PHONY: k8s-deploy-other
k8s-deploy-other:
	kustomize build k8s/other | kubectl apply -f -
	kustomize build k8s/csr | kubectl apply -f -
	@echo Waiting for cert creation ...
	@sleep 15
	kubectl certificate approve $(WEBHOOK_SERVICE).$(NAMESPACE)

.PHONY: k8s-deploy-csr
k8s-deploy-csr:
	kustomize build k8s/csr | kubectl apply -f -
	@echo Waiting for cert creation ...
	@sleep 15
	kubectl certificate approve $(WEBHOOK_SERVICE).$(NAMESPACE)

.PHONY: k8s-deploy-deployment
k8s-deploy-deployment:
	(cd k8s/deployment && \
	kustomize edit set image CONTAINER_IMAGE=$(CONTAINER_IMAGE))
	kustomize build k8s/deployment | kubectl apply -f -

k8s-delete-all:
	kustomize build k8s/other | kubectl delete --ignore-not-found=true -f  - 
	kustomize build k8s/csr | kubectl delete --ignore-not-found=true -f  - 
	kustomize build k8s/deployment | kubectl delete --ignore-not-found=true -f  - 
	kubectl delete --ignore-not-found=true csr $(WEBHOOK_SERVICE).$(NAMESPACE)
	kubectl delete --ignore-not-found=true secret hello-tls-secret


test:
	cd webhook && go test ./...