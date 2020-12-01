# Mirror docker repositories to another registry

This tool scans for matched tags on a given repository and synchronizes it with own registry.

Pull the docker image with the compiled binary

```
docker pull bdudelsack/docker-registry-mirror
```

Create `config.yaml` with configuration

```yaml
repositories:
  - source: library/nginx
    destination: registry.example.com/nginx
    matches:
      - ^1\.11(\.[0-9]+)?(-alpine)?$
  - source: library/php
    destination: registry.example.com/php
    matches:
      - ^7\.1\.2
      - ^5\.6(-fpm)?$
  - source: library/debian
    destination: registry.example.com/debian
    matches:
      - ^jessie(-slim)?$
      - ^latest$
  - source: quay.io/coreos/hyperkube
    destination: registry.example.com/hyperkube
    matches:
      - ^v1\.(5|6)\.[0-9]+_coreos\.[0-9]+$
auth:
  registry.example.com:
    username: myusername
    password: mypassword
  registry.hub.docker.com:
    username: myusername
    password: mypassword
  quay.io:
    username: ${MY_USERNAME}
    password: ${MY_PASSWORD}
```

Run the docker container with mounted config file

```
docker run -v $PWD/config.yaml:/config.yaml bdudelsack/docker-registry-mirror
```

Specify credentials with environment variables

```
docker run -v $PWD/config.yaml:/config.yaml -e MY_USERNAME=username -e MY_PASSWORD=password bdudelsack/docker-registry-mirror
```