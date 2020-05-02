---
title: "Install"
menu: true
weight: 3
---

You can install the pre-compiled binary (in several different ways), use Docker or compile from source.

Here are the steps for each of them:

## Install the pre-compiled binary

Download the pre-compiled binaries from the [releases page](https://github.com/shellbear/web-watcher/releases) and
copy to the desired location.

## Running with Docker

You can also run it within a Docker container. Here as follows an example command:

```sh
export DISCORD_TOKEN=XXXXXXXXX....
```

```sh
docker run -d                                               \
    -e DISCORD_TOKEN                                        \
    -v web-watcher-data:/app                                \
    --name web-watcher                                      \
    docker.pkg.github.com/shellbear/web-watcher/web-watcher
```

The container is based on latest Go docker image.

## Compiling from source

If you feel adventurous you can compile the code from source:

```sh
git clone https://github.com/shellbear/web-watcher.git
cd web-watcher

# get dependencies using go modules (needs go 1.11+)
go get ./...

# build
go build -o web-watcher .

# check it works
./web-watcher
``` 