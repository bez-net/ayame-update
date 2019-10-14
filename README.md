# WebRTC Signaling Server Ayame Update
    upgraded by Stoney Kang(CSO), sikang99@gmail.com in TeamGRIT


## History
- 2019/10/13 add endpoints of `/admin` and `/event`
- 2019/10/12 add local db handling (19.03.01)
- 2019/10/12 add Common for structs such as client, room, hub
- 2019/10/11 /admin endpoint added to check the server info
- 2019/10/10 not use logrus because of readability
- 2019/10/10 add `MaxSessions` in the server configuration
- 2019/10/10 removed `ws_handler.go`
- 2019/10/09 update to support plain(3000) and secure(3443) server at the same time
    - certs/ directory includes *.pem files which are only for localhost
- 2019/07/08 translate jp into kr for files in the foler `/doc`
- 2019/07/08 translate Japanese(jp) into Korean(kr) for files in the foler `/doc`
- 2019/07/08 update slightly `Makefile`, `README.md`
    - Refer the base example [webrtc/apprtc](https://github.com/webrtc/apprtc) : video chatting demo app
- 2019/07/07 forked from [OpenAyame/ayame](https://github.com/OpenAyame/ayame)


## Information
- [hakobera/serverless-webrtc-signaling-server](https://github.com/hakobera/serverless-webrtc-signaling-server) - Serverless WebRTC Signaling Server only works for WebRTC P2P
