# Tsnet Reverse Proxy for Docker Compose

## What:
This is a docker container that allows you to easily reverse-proxy other services as a custom tsnet device.

## Wait, what?

If you have a docker service running in a container. This allows you to expose that container as a custom tailscale device. It will get its own tailscale-ip-address, and custom hostname. It will also optionally support HTTPS, if that is setup on your tailnet.


## Why:
If you have multiple services on a single device (in docker-compose) this allows you to break them out into separate tailscale devices.
for example, on my tailscale network, my raspberry pi is running many services. However, they are all accessed as if they are separate devices. "pihole" is one device, "photoprism" is another, etc.


## How:
An example `docker-compose.yml`. In this case, putting `pihole`'s web interface behind a reverse proxy.
```yml
version: '3'
services:
  pihole:
    container_name: pihole
    image: pihole/pihole:latest
    ports:
      - "53:53/tcp"
      - "53:53/udp"
      # no need to expose port 80. it is exposed through tailscale
    # .. other pihole config excluded for brevity
  tsnetproxy:
    container_name: pihole_proxy
    image: something
    environment:
      TS_AUTHKEY: $TS_AUTHKEY
      TSNET_HOSTNAME: pihole
      TSNET_HOST: pihole:80 #this refers to the "pihole" container above. docker containers can access themselves by name when on same network.
      
```


## Building
```bash
docker build -t tsnet-composable .
```