# acme-fixer

- [Binaries](https://github.com/containous/acme-fixer/releases)
- [Docker image](https://hub.docker.com/r/containous/acme-fixer)
- [Documentation](./docs/acme-fixer.md)

## Examples

### Traefik v1

With binaries:

```bash
# dry run mode
acme-fixer -i ./my/path/acme.json -d

# without dry run
acme-fixer -i ./my/path/acme.json
```

With Docker:

```bash
# dry run mode
docker run -v $PWD/letsencrypt:/letsencrypt containous/acme-fixer:v0.1.1 -i /letsencrypt/acme.json -d

# without dry run
docker run -v $PWD/letsencrypt:/letsencrypt containous/acme-fixer:v0.1.1 -i /letsencrypt/acme.json
```

### Traefik v2

With binaries:

```bash
# dry run mode
acme-fixer -i ./my/path/acme.json -d --v2

# without dry run
acme-fixer -i ./my/path/acme.json --v2
```

With Docker:

```bash
# dry run mode
docker run -v $PWD/letsencrypt:/letsencrypt containous/acme-fixer:v0.1.1 -i /letsencrypt/acme.json -d --v2

# without dry run
docker run -v $PWD/letsencrypt:/letsencrypt containous/acme-fixer:v0.1.1 -i /letsencrypt/acme.json --v2
```
