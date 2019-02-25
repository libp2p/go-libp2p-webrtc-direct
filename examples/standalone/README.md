standalone
===

This example shows the transport being used on its own. It also shows off compatibility with [js-libp2p-webrtc-direct](https://github.com/libp2p/js-libp2p-webrtc-direct).

## Go

### Install dependencies
**TODO**: Check the root readme

### Listener
```sh
go run main.go -listen
```
*Output*
```
[listener] Listening
[listener] Got connection
[listener] Got stream
[listener] Received:
hey, how is it going. I am the dialer
Failed to accept data channel: The association is closed
```
The last line is harmless warning printed by the pions/webrtc library.
### Dialer
```sh
go run main.go
```
*Output*
```
Warning: Certificate not checked
[dialer] Opened connection
[dialer] Opened stream
Failed to push SCTP packet: Failed sending reply: dtls: conn is closed
Warning: mux: no endpoint for packet starting with 23
Failed to push SCTP packet: Failed sending reply: dtls: conn is closed
Warning: mux: no endpoint for packet starting with 21
Failed to accept data channel: The association is closed
```
The warnings printed by the pions/webrtc library are harmless.

## Javascript
The equivalent javascript example is also provided. It can be used as follows:

### Install dependencies
```sh
npm install
```

### Listener
```sh
node index.js --listen
```
*Output*
```
[listener] Listening
[listener] Got connection
[listener] Got stream
[listener] Received:
hey, how is it going. I am the dialer
```
### Dialer
```sh
node index.js
```
*Output*
```
[dialer] Opened connection
[dialer] Opened stream
```