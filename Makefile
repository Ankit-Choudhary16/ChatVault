.PHONY: build run clean db-up killport
MY_APP = chatvault
PORT = 8080
GO = go
.DEFAULT_GOAL = build
#local deployment
killport:
	@lsof -ti :$(PORT) | xargs -r kill -9|true
build:
	@mkdir -p bin 
	$(GO) build -o ./bin/$(MY_APP) ./cmd/server

run:build killport
	cp .env.example .env
	$(GO) run ./cmd/server


db-up:
	@echo "Starting local PostgreSQL..."
	brew services start postgresql@14

clean: 
	rm -rf ./bin



#k8s deployment in kind cluster
create-cluster:
	kind create cluster --name chatvault-cluster --config deployments/kind-config.yaml
delete-cluster:
	kind delete cluster --name chatvault-cluster

load-image:
	kind load docker-image chatvaultcopy2-app:latest --name chatvault-cluster

ingress-controller:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

#always apply ingress after controller is up
deploy:
	kubectl apply -k deployments/k8s

chatvault.local:
	echo "127.0.0.1 chatvault.local" | sudo tee -a /etc/hosts


delete-resources: 
	kubectl delete all --all -n chatvault


delete-ingress:
	kubectl delete ingress chatvault-ingress -n chatvault

delete-ingress-controller:
	kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml



delete-namespace:
	kubectl delete namespace chatvault

#helm 

deploy-helm:
	helm install chatvault-release ./chatvault \
  	-n chatvault \
  	--create-namespace

upgrade-helm:
	helm upgrade chatvault-release ./chatvault \
  	-n chatvault --values ./chatvault/values.yaml
delete-helm:
	helm uninstall chatvault-release -n chatvault