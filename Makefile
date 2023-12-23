help:
	@echo "Use 'make <ROOT>' where <ROOT> is one of"
	@echo "  minikube:            to download and install minikube on linux X86"
	@echo "  minikube-start:      to start the cluster"
	@echo "  minikube-stop:       to halt the cluster"


minikube:
	curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
	sudo install minikube-linux-amd64 /usr/local/bin/minikube

minikube-start:
	minikube start

minikube-stop:
	minikube stop

