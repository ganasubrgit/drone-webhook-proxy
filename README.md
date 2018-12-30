# drone-webhook-proxy
Proxy and agent to convert passive trigger of Drone CI to active polling trigger.

## Motivation and Design

I love [Drone CI](https://drone.io) as a modern continuous integration server for its first class docker support and simple concepts.

However, one hitch with Drone is that it is only triggered by repository web hooks and offers no way to periodically poll repositories for changes.

This is a deal breaker for deployments in private network environments (i.e. behind firewall, home network).

I made this utility to act as a web hook proxy for Drone in order to convert passive trigger into active polling.

On the cloud side, `hook proxy` should receive web hook calls instead of Drone server. And on the deployment side, `hook agent` will be configured to periodically poll `hook proxy` for events and deliver that to the local Drone instance.

## Basic Usage 

**On the cloud side**, spin up a VM. For instance, start a docker droplet on DigitalOcean.

```bash
cat >> docker-compose.yml << EOF
version: "3"

services:
  redisdb:
    image: "bitnami/redis:latest"
    environment:
    - "ALLOW_EMPTY_PASSWORD=yes"
    volumes:
    - "redis-data:/bitnami/redis/data"
  proxy:
    image: "davidiamyou/drone-webhook-proxy:LATEST"
    links:
      - redisdb
    depends_on:
      - redisdb
    expose:
      - "8080"
    ports:
      - "35027:8080"
    entrypoint: hook proxy -r redisdb:6379 -x 500

volumes:
  redis-data:
EOF

docker-compose up -d

# Don't forget to open up ports
sudo ufw allow 35027/tcp
```

**On local side**, start `hook agent`.

```bash
# Assuming Drone is available at http://192.168.100.100:8080
# Assuming remote hook proxy is available at http://my-hook-proxy
docker run -d davidiamyou/drone-webhook-proxy:LATEST \
    agent \
    -p http://my-hook-proxy/pop \
    -d http://192.168.100.100:8080/hook
```

**Lastly, change the web hook address on CVS system to `http://my-hook-proxy/push`**, assuming `hook proxy` is available at `http://my-hook-proxy`.

