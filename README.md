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
