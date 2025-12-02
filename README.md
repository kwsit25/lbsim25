# lbsim25
loadbalancing simulation in k8s

## Install

1. Install k8s (enable via Docker Desktop)
1. Install go [Download and install](https://go.dev/doc/install)
1. Build docker container of services
	1. Install go dependencies with `go mod download`
	1. Build go webclient & go webserver
		1. go into the folders (`services/webserver` & `services/webclient`)
		1. execute `make buildcontainerapp`
	1. Build Docker image for all (client, server, nginx)
		1. go into the folders (`services/webserver`, `services/webclient` & `services/nginx`)
		1. execute `make docker`


## TODO

1. Install & setup monitoring
1. Apply server k8s deployement(s) (see example at `k8s/sample`)
1. Create k8s deplyoment for service to server of type `Loadbalander`
1. Open browser on defined ports of service e.g. `http://localhost:8080/api/load?source=manual`
1. Check if metrics are available
1. Build grafana dashboards with existing metrics
1. Concider loadbalancer config for nginx
1. Setup loadbalancing with nginx and wanted strategies
