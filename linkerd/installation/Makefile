install_linkerd_local:
	curl -sL https://run.linkerd.io/install | sh

check_environment_local:
	linkerd check --pre

install_linkerd_cluster:
	linkerd install | kubectl apply -f -

check_environment_cluster:
	linkerd check

start_dashboard:
	linkerd dashboard &
