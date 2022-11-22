# ⚠️⚠️⚠️⚠️⚠️⚠️
**Status:**

[Archived](https://github.com/libp2p/github-mgmt/pull/74) and not maintained

**Alternatives:**

WebRTC Browser-to-Server is being implemented in go-libp2p here https://github.com/libp2p/specs/pull/412 per the specification: https://github.com/libp2p/specs/pull/412

WebRTC Browser-to-Browser is being tracked here: https://github.com/libp2p/specs/issues/475

**Questions:**

Please direct any questions about the specification to: https://github.com/libp2p/specs/issues

Please direct any questions about the go-libp2p WebRTC implementation to: https://github.com/libp2p/go-libp2p/issues
# ⚠️⚠️⚠️⚠️⚠️⚠️

# go-libp2p-webrtc-direct

[![](https://img.shields.io/badge/project-libp2p-yellow.svg?style=flat-square)](http://github.com/libp2p/libp2p)
[![](https://img.shields.io/badge/freenode-%23libp2p-yellow.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23libp2p)
[![GoDoc](https://godoc.org/github.com/libp2p/go-libp2p-webrtc-direct?status.svg)](https://godoc.org/github.com/libp2p/go-libp2p-webrtc-direct)
[![Coverage
Status](https://coveralls.io/repos/github/libp2p/go-libp2p-webrtc-direct/badge.svg?branch=master)](https://coveralls.io/github/libp2p/go-libp2p-webrtc-direct?branch=master)
[![Build
Status](https://travis-ci.org/libp2p/go-libp2p-webrtc-direct.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-webrtc-direct)

> A transport that enables browser-to-server, and server-to-server, direct
> communication over WebRTC without requiring signalling servers. This is the
> Go counterpart to
> [js-libp2p-webrtc-direct](https://github.com/libp2p/js-libp2p-webrtc-direct).

Lead maintainer: [@backkem](https://github.com/backkem)

Special thanks to [@pion](https://github.com/pion) for their fantastic
[WebRTC Go library](https://github.com/pion/webrtc), which made this
libp2p transport possible.

## Install

This package supports gomod builds.

```sh
go get github.com/libp2p/go-libp2p-webrtc-direct
```

## Usage

Check out the
[GoDocs](https://godoc.org/github.com/libp2p/go-libp2p-webrtc-direct).

## Examples

Check the [examples](./examples) folder for usage and integration examples.

## Contribute

Feel free to join in. All welcome. Open an
[issue](https://github.com/libp2p/go-libp2p-webrtc-direct/issues) or send a
PR.

This repository falls under the IPFS [Code of
Conduct](https://github.com/ipfs/community/blob/master/code-of-conduct.md).

## License
MIT