# Authentication

All API commands require DCI RemoteCI credentials. You can obtain credentials from the [DCI dashboard](https://www.distributed-ci.io/).

## Option 1: Config File

```bash
go-dci config set --accesskey <your-access-key> --secretkey <your-secret-key>
```

This creates a `.go-dci-config.yaml` file in the current directory.

```
Usage:
  dci config set [flags]

Flags:
  -a, --accesskey string   The access key to set in the configuration.
  -h, --help               help for set
  -s, --secretkey string   The secret key to set in the configuration.
```

## Option 2: Environment Variables

You can also set credentials via environment variables (useful for CI/CD):

```bash
export GO_DCI_ACCESSKEY=<your-access-key>
export GO_DCI_SECRETKEY=<your-secret-key>
```

| Variable | Description |
|----------|-------------|
| `GO_DCI_ACCESSKEY` | Your DCI client ID / access key |
| `GO_DCI_SECRETKEY` | Your DCI API secret key |

Environment variables take precedence over values in the config file.

## Additional Configuration

### OCP Version Tracking

The `ocpcount` command tracks specific OCP versions by default. You can customize which versions to track using the `OCP_VERSIONS_TO_TRACK` environment variable:

```bash
export OCP_VERSIONS_TO_TRACK="4.15,4.16,4.17"
```

If not set, the following versions are tracked by default: 4.12, 4.13, 4.14, 4.15, 4.16, 4.17, 4.18, 4.19, 4.20.

| Variable | Description | Default |
|----------|-------------|---------|
| `OCP_VERSIONS_TO_TRACK` | Comma-separated list of OCP versions to track in `ocpcount` command | 4.12, 4.13, 4.14, 4.15, 4.16, 4.17, 4.18, 4.19, 4.20 |

---

Next: [CLI Reference](cli-reference.md)
