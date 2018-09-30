const WebRTCDirect = require('libp2p-webrtc-direct')
const multiaddr = require('multiaddr')
const pull = require('pull-stream')

const mh = multiaddr('/ip4/127.0.0.1/tcp/9090/http/p2p-webrtc-direct')

const direct = new WebRTCDirect()

const listener = direct.createListener({config:{}},(socket) => {
  console.log('new connection opened')
  pull(
    pull.values(['hello']),
    socket
  )
})

listener.listen(mh, () => {
  console.log('listening')
})