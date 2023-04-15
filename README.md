# Tsnet Reverse Proxy for Docker Compose

## What:
This is a docker container that allows you to easily reverse-proxy other services as a custom tsnet device.

## Wait, what?

If you have a docker service running in a container. This allows you to expose that container as a custom tailscale device. It will get its own tailscale-ip-address, and custom hostname. It will also optionally support HTTPS, if that is setup on your tailnet.


## Why:
If you have multiple services on a single device (in docker-compose) this allows you to break them out into separate tailscale devices.
for example, on my tailscale network, my raspberry pi is running many services. However, they are all accessed as if they are separate devices. "pihole" is one device, "photoprism" is another, etc.


## How:
An example `docker-compose.yml`. In this case, putting `vaultwarden`'s web interface behind a reverse proxy.
This is a great use case for it, making vaultwarden completely private, and only available over a tailnet. While also completely containerized.
No ports to bitwarden are exposed, even to localhost. Its only exposed to the other container, which then only exposes it to the tailnet over https.
```yml
version: '2.1'
services:
  vaultwarden:
    image: vaultwarden/server:latest
    container_name: vaultwarden
    volumes:
      - vaultwarden-data:/data
    restart: unless-stopped
  vwexperimentproxy:
    image: tsnet-composable-stable
    container_name: tsnet-composable-stable
    environment:
      - TS_AUTHKEY=tskey-auth-SOMEKEYHERE #get a key from the tailscale site under settings.
      - TSNET_CUSTOM_HOSTNAME=vaultwarden #Give it a custom name if you want. in this case, its then available at https://vaultwarden.ts-net-name.ts.net
      - TSNET_PROXY_TO_URL=http://vaultwarden  #this URL refers to the other container's container_name above
    volumes:
      - tailscale-data:/var/lib/tailscale
    restart: unless-stopped
volumes:
  vaultwarden-data:
  tailscale-data:

```


## Building
```bash
docker build -t tsnet-composable-stable .
```