# iscn-ipld

This repository is a go-ipfs IPLD plugin for [International Standard Content Number (ISCN)](https://github.com/likecoin/iscn-specs).

## Building and Installing
You can build the  plugin by running `make build`. You can then install it into your local IPFS repo by running `make install`.

Plugins need to be built against the correct version of go-ipfs. This package generally tracks the latest go-ipfs release but if you need to build against a different version, please set the `IPFS_VERSION` environment variable.

You can set `IPFS_VERSION` to:

* `vX.Y.Z` to build against that version of IPFS.
* `$commit` or `$branch` to build against a specific go-ipfs commit or branch.
* `/absolute/path/to/source` to build against a specific go-ipfs checkout.

To update the go-ipfs, run:

```
> make go.mod IPFS_VERSION=version
```
