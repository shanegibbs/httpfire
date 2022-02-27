# httpfire 🔥

Accurately test your web infrastructure. Understand your limits and failure scenarios. Test and validate scaling strategies.

[Roadmap](https://github.com/shanegibbs/httpfire/projects/1)

## Development

Requires:
- `go >= 1.17.6`
- `make`
- `docker`
- `docker-compose`

Run in local mode for development

```shell
$ make local
```

## Docker

Build and run a local cluster with docker compose

```shell
$ make up
```

Start a test plan

```shell
$ make post-start
```

Stop current test plan

```shell
$ make post-stop
```

Cleanup (remove docker containers etc)

```shell
$ make clean
```
