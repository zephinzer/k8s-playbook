start_minikube:
	@$(MAKE) start_minikube_mac

_start_minikube:
	minikube \
		--vm-driver=${VM_DRIVER} \
		--memory=2000 \
		--disk-size=5g \
		start

start_minikube_mac:
	@$(MAKE) _start_minikube VM_DRIVER=hyperkit

start_minikube_linux:
	@$(MAKE) _start_minikube VM_DRIVER=kvm2

use_minikube_docker:
	eval $$(minikube docker-env)

stop_minikube:
	minikube stop
