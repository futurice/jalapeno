---
sidebar_position: 2
slug: /installation/
title: Installation
---

# Installation

## From the Binary Releases

Every release of Jalapeno provides binary releases for a variety of OSes. These binary versions can be manually downloaded and installed.

1. Download your [desired version](https://github.com/futurice/jalapeno/releases)
1. Make the binary executable (`chmod +x jalapeno-linux-amd64`)
1. Rename and move the binary to proper location (`mv jalapeno-linux-amd64 /usr/local/bin/jalapeno`)

## Build From Source

First, make sure you have [Go](https://go.dev/doc/install) and [Task](https://taskfile.dev/installation) installed.

Then you can compile the binary with following commands:

```bash
git clone https://github.com/futurice/jalapeno.git
cd jalapeno
task build
```

After this the binary is available on path `./bin/jalapeno`.
