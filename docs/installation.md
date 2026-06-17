# Installation

## Prebuilt Binary

Download the latest release for your platform from [GitHub Releases](https://github.com/sebrandon1/go-dci/releases):

```bash
# Linux (amd64)
curl -sL https://github.com/sebrandon1/go-dci/releases/latest/download/go-dci_$(curl -sL https://api.github.com/repos/sebrandon1/go-dci/releases/latest | grep tag_name | cut -d '"' -f4)_linux_amd64.tar.gz | tar xz
sudo mv go-dci /usr/local/bin/

# macOS (Apple Silicon)
curl -sL https://github.com/sebrandon1/go-dci/releases/latest/download/go-dci_$(curl -sL https://api.github.com/repos/sebrandon1/go-dci/releases/latest | grep tag_name | cut -d '"' -f4)_darwin_arm64.tar.gz | tar xz
sudo mv go-dci /usr/local/bin/
```

## Container Image

```bash
# Using podman
podman run --rm \
  -e GO_DCI_ACCESSKEY="$GO_DCI_ACCESSKEY" \
  -e GO_DCI_SECRETKEY="$GO_DCI_SECRETKEY" \
  quay.io/bapalm/go-dci identity

# Using docker
docker run --rm \
  -e GO_DCI_ACCESSKEY="$GO_DCI_ACCESSKEY" \
  -e GO_DCI_SECRETKEY="$GO_DCI_SECRETKEY" \
  quay.io/bapalm/go-dci identity
```

## Go Install

```bash
go install github.com/sebrandon1/go-dci@latest
```

## Build from Source

```bash
git clone https://github.com/sebrandon1/go-dci.git
cd go-dci
make build
```

---

Next: [Authentication](authentication.md)
