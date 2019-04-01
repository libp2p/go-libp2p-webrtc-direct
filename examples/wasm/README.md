wasm
===

This example demos WASM support.

The example is supposed to be ran against one of the `standalone` examples in `listening` mode. Note that the example only supports the `Dial` side since `Listen` uses an HTTP server which isn't available in a browser.

## Usage

### Install dependencies
**TODO**: Check the root readme

### Run

```sh
GOOS=js GOARCH=wasm go build -o main.wasm
```
Next, refer to the [Go documentation](https://github.com/golang/go/wiki/WebAssembly#getting-started) for how to run a wasm file.

*Output in the browser console*
```
[dialer] Opened connection wasm_exec.js:47:6
[dialer] Opened stream
```
*Output in the listener*
```
[listener] Got stream
[listener] Received:
hey, how is it going. I am the dialer
```