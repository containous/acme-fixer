# acme-fixer

[Documentation](./docs/acme-fixer.md)

## Examples

### Traefik v1

```bash
# dry run mode
acme-fixer -i ./my/path/acme.json -d

# without dry run
acme-fixer -i ./my/path/acme.json
```

### Traefik v2

```bash
# dry run mode
acme-fixer -i ./my/path/acme.json -d --v2

# without dry run
acme-fixer -i ./my/path/acme.json --v2
```
