## WebRTC Signaling Server Ayame Update
    updated by Stoney Kang(CSO), sikang99@gmail.com in TeamGRIT

### Usage
1. build the server and run it
```
$ make build
$ make run # or ./ayame
```
2. connect to server to get web page on your browser
```
$ open http://localhost/static # use xdg-open in Linux instead of open of MacOS
```


### API endpoints
- `/static` : file server for web pages
- `/signal` : webrtc signaling
- `/admin` : status monitoring
- `/event` : handle push/pull event with sub points of `/event/send|/event/recv`
- `/upload` : file upload for sharing


### History
- 2019/10/16
    - remove .circleci directory because of no more use
    - add functions to send panic stack message to slack
- 2019/10/15 add `/upload` to support file upload
    - move `/sample` into `/static` as an endpoint for file server
    - add self health checker for running
- 2019/10/14 rename `/signaling` into `/signal`
    - handle SSE events
- 2019/10/13 add endpoints of `/admin` and `/event`
- 2019/10/12 add local db handling (19.03.01)
- 2019/10/12 add Common for structs such as client, room, hub
- 2019/10/11 `/admin` endpoint added to check the server info
- 2019/10/10 not use logrus because of readability, some day later return to use
- 2019/10/10 add `MaxSessions` in the server configuration
- 2019/10/10 removed `ws_handler.go`
- 2019/10/09 update to support plain(3000) and secure(3443) server at the same time
    - certs/ directory includes *.pem files which are only for localhost
- 2019/07/08 translate jp into kr for files in the foler `/doc`
- 2019/07/08 translate Japanese(jp) into Korean(kr) for files in the foler `/doc`
- 2019/07/08 update slightly `Makefile`, `README.md`
    - Refer the base example [webrtc/apprtc](https://github.com/webrtc/apprtc) : video chatting demo app
- 2019/07/07 forked from [OpenAyame/ayame](https://github.com/OpenAyame/ayame)


### Information
- [Signaling and video calling](https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API/Signaling_and_video_calling)
- [hakobera/serverless-webrtc-signaling-server](https://github.com/hakobera/serverless-webrtc-signaling-server) - Serverless WebRTC Signaling Server only works for WebRTC P2P
- [FiloSottile/mkcert](https://github.com/FiloSottile/mkcert) - A simple zero-config tool to make locally trusted development certificates with any names you'd like. https://mkcert.dev


### Reference
- [How to format current time using a yyyyMMddHHmmss format?](https://stackoverflow.com/questions/20234104/how-to-format-current-time-using-a-yyyymmddhhmmss-format)
- [Uploading Files in Go - Tutorial](https://tutorialedge.net/golang/go-file-upload-tutorial/)

