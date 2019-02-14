const flags = require('flags')
const WebRTCDirect = require('libp2p-webrtc-direct')
const multiaddr = require('multiaddr')
const mplex = require('libp2p-mplex')
const pull = require('pull-stream')

const listenFlag = 'listen'
flags.defineBoolean(listenFlag, false, 'Listen for incoming connections.')
flags.parse()
const listening = flags.get(listenFlag)

const maddr = multiaddr('/ip4/127.0.0.1/tcp/9090/http/p2p-webrtc-direct')

const direct = new WebRTCDirect()

if (listening) {
  const listener = direct.createListener({ config: {} }, (conn) => {
    console.log('[listener] Got connection')

    const muxer = mplex.listener(conn)

    muxer.on('stream', (stream) => {
      console.log('[listener] Got stream')
      pull(
      stream,
      pull.drain((data) => {
        console.log('[listener] Received:')
        console.log(data.toString())
      })
      )
    })
  })

  listener.listen(maddr, () => {
    console.log('[listener] Listening')
  })
} else {
  direct.dial(maddr, { config: {} }, (err, conn) => {
    if (err) {
      console.log(`[dialer] Failed to open connection: ${err}`)
    }
    console.log('[dialer] Opened connection')

    const muxer = mplex.dialer(conn)
    const stream = muxer.newStream((err) => {
      console.log('[dialer] Opened stream')
      if (err) throw err
    })

    pull(
      pull.values(['hey, how is it going. I am the dialer']),
      stream
    )
  })
}
