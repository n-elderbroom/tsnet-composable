# Tsnet Reverse Proxy for Docker Compose

## What:
This is a simple docker container that allows you to easily reverse-proxy other containers as a custom tsnet device.

Specifically, this is designed to wrap containers that offer an http server, and expose those over tailscale, with HTTPS support.
It does this while keeping the http server private. The server is *only* exposed over tailscale. (you can make it public via tailscale funnel!)

I think in theory, this might be possible with the official tailscale container? but i couldn't quite figure it out. They have a guide that I think does this for kubernetes? but i've just got a small docker-compose setup. :shrug: I tried `tailscale serve` on the official container, but it is only designed to serve *local* resources, so tailscale has to be installed in the same container as the service. I didn't want to rebuild the containers for every service I want to host. The official `tsnet` package on the other hand is quite simple, and `go` has built in reverse-http-proxy stuff. 

Just recently, some docker mods for linuxserver.io containers got released, which really does all the hard work for you. If the service you want has a linuxserver.io container, you should use that method instead of this one. It is a much cleaner solution in my opinion.



## Why:
If you have multiple services on a single device (in docker-compose) this allows you to break them out into separate tailscale devices.
for example, on my tailscale network, my raspberry pi is running many services. However, using this, I can make them appear as separate tailscale devices.
This means that they are all on separate subdomains. This makes some services happier (since every service gets the root of its own subdomain), but also just felt more organized to me. You could also share out -and control access to- specific services very easily this way, using tailscale's sharing features and ACLs.  


## Example:
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
  vaultwardenproxy:
    image: ghcr.io/n-elderbroom/tsnet-composable:main
    container_name: tsnet-composable-stable
    environment:
      - TS_AUTHKEY=tskey-auth-SOMEKEYHERE #get a key from the tailscale site under settings.
      - TSNET_CUSTOM_HOSTNAME=vaultwarden #Give it a custom name if you want. in this case, its then available at https://vaultwarden.ts-net-name.ts.net
      - TSNET_PROXY_TO_URL=http://vaultwarden  #this URL refers to the other container's container_name above
      # - TSNET_ENABLE_FUNNEL=true #uncomment this to make the container publicly available over tailscale funnel. 
    volumes:
      - tailscale-data:/root/.config/tsnet-app
    restart: unless-stopped
volumes:
  vaultwarden-data:
  tailscale-data:

```




## Building
```bash
docker build -t tsnet-composable-stable .
```