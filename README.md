Package ``go-libp2p-webrtc-direct`` is a Golang version of the webrtc-direct libp2p transport.

Please refer to ``pions/webrtc`` for additional installation instructions.
The package also requires the following forks to be checked out under their original package name:
- ``backkem/go-multiaddr``
- ``backkem/mafmt``

The transport passes the ``SubtestStress1Conn1Stream1Msg`` test case but there is a long list of known limitations. Therefore, please don't rely on this package. It only serves as a proof of concept and as an experiment to gather some experience building tools on top of ``pions/webrtc``.