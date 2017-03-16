# Mirror docker repositories to another registry

[Configuration example](./config.yaml.example)

```
docker pull bdudelsack/docker-registry-mirror
docker run -v $PWD/config.yaml:/config.yaml -v /etc/ssl/certs:/etc/ssl/certs bdudelsack/docker-registry-mirror
```