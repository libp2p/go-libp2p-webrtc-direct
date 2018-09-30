
const WebRTCDirect = require('libp2p-webrtc-direct')
const multiaddr = require('multiaddr')
const pull = require('pull-stream')

const mh = multiaddr('/ip4/127.0.0.1/tcp/9090/http/p2p-webrtc-direct')

const direct = new WebRTCDirect()


direct.dial(mh, (err, conn) => {
      if(err){
        console.log(`Error: ${err}`)
      }

        console.log(`dial success`)
  pull(
    conn,
    pull.collect((err, values) => {
      if (!err) {
        console.log(`Value: ${values.toString()}`)
      } else {
        console.log(`Error: ${err}`)
      }
    }),
  )
  })