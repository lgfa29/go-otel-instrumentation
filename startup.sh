#!/bin/bash

# Get the IP of the host machine (host.docker.internal)
IP=$(echo $(ping -c 1 host.docker.internal | gawk -F'[()]' '/PING/{print $2}'))

# Add the host IP and mappings to the OTel Collector and Traefik to /etc/hosts
echo "$IP host.docker.internal traefik.localhost otel-collector-http.localhost" >> /etc/hosts

# Start the go program
printf "Now starting $1.go...\n"
go run $1.go