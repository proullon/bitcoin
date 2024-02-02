## Bitcoin Go implementation

Naive implementation of bitcoin in Go, using [bitcoin whitepaper](https://bitcoin.org/bitcoin.pdf).

## Build and run

There is to binaries, a CLI and a network node. To build them:

```sh
$ make install
$ bitcoin-cli help
$ bitcoin-node help
```

Generate a key:

```sh
$ bitcoin-cli generate-key secret.key
```

Then start the node:

```sh
$ bitcoin-node run --key secret.key --workdir /tmp/
```
