---
sidebar_position: 2
slug: /installation/
title: Installation
---

# Installation

## From the Binary Releases

Cross-platform binaries are provided with each release of Jalapeno. These can manually be
downloaded and installed from [GitHub releases](https://github.com/futurice/jalapeno/releases/).

In short the process is:

1. Download [the latest version of jalapeno](https://github.com/futurice/jalapeno/releases/latest)
for your platform
2. Make the binary executable
3. Rename and move the binary to proper location

For example on MacOS (running on Apple Silicon) this can be done with:

```bash
curl -L https://github.com/futurice/jalapeno/releases/latest/download/jalapeno-darwin-arm64 -o jalapeno
chmod +x jalapeno
mv jalapeno /usr/local/bin/jalapeno
```

## Use via Docker

It is possible to use jalapeno via Docker without installing it locally. Jalapeno image is available
from the GitHub Container Registry, and thus the following command on a *nix system is equivalent
to running `jalapeno` locally:

MacOS and Linux:

```bash
docker run -it --rm -v $(pwd):/workdir ghcr.io/futurice/jalapeno:v0.1.4
```

Windows Command Line:

```batch
docker run -it --rm -v %cd%:/workdir ghcr.io/futurice/jalapeno:v0.1.4
```

PowerShell:

```powershell
docker run -it --rm -v ${PWD}:/workdir ghcr.io/futurice/jalapeno:v0.1.4
```

## Build From Source

First, make sure you have [Go](https://go.dev/doc/install) and
[Task](https://taskfile.dev/installation) installed.

Then you can compile the binary with following commands:

```bash
git clone https://github.com/futurice/jalapeno.git
cd jalapeno
task build
```

After this the binary is available on path `./bin/jalapeno`.
