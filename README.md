# DNS IP Updater

This program polls the host's ip public address and infroms a third-party server to update the DNS records with the host's new IP.

The client is the 'host' in this project. It is the computer that runs the services and its where the DNS records point to. The server is the third-party computer that is in charge of updating the DNS records for the client because the client's new IP is not whitelisted by namecheap's API. 

## Getting Started

Create the file 'config.go' in both the client and server 

server/config.go
```
package main
 
const namecheapAPIUser = ""
const namecheapAPIToken = ""
const namecheapUserName = ""
const namecheapSLD = "yahoo"
const namecheapTLD = "com"
```


client/config.go
```
package main
 
const twilioSid = ""
const twilioToken = ""
const twilioFrom = "" // format: +17778889999
const twilioTo = ""

const serverURL = ""
```

## Testing Locally

To run the server
```
KEY="default" go run server/*.go
```

## Build Docker

The Client
```
docker build -t registry.tunerinc.com/dns-ip-updater:client ./client/
```
The Server
```
docker build -t registry.tunerinc.com/dns-ip-updater:server ./server/
```

## Running The Program

### Locally

Run the Client 
```
docker run -d --restart=unless-stopped /
--name=dns-ip-updater-client \
-e KEY=default -e USERNAME=username -e PASSWORD=password \
registry.tunerinc.com/dns-ip-updater:client
```

Run the Server  
```
docker run -d --restart=unless-stopped /
--name=dns-ip-updater-server -e KEY=default \
registry.tunerinc.com/dns-ip-updater:server
```

### Production (On host and VM) 

Fill in the docker-compose template for both the client and server
```
cp docker-compose.template.yaml docker-compose.yaml 
```

Fill in traefik config
```
cp traefik.template.toml traefik.toml
```

To encrypt the login information for basic-auth use this command after filling out your username and password.
```
echo $(htpasswd -nb <AUTH-USER> <AUTH-PASS>) | sed -e s/\\$/\\$\\$/g
```

Run the client with the docker-compose file
```
docker-compose up -d
```

Run the Server with docker-compose file on a GCE VM(shared cpu + .6Gb RAM)
```
docker network create web && \
touch ./dns-ip-updater/server/acme.json && chmod 600 ./dns-ip-updater/server/acme.json && \
docker-compose up -d
```
