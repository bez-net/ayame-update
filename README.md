## WebRTC Signaling Server Ayame Update
updated by Stoney Kang(CSO), sikang99@gmail.com in [TeamGRIT.kr](https://teamgrit.kr)

### Usage
1. build the server and run it
```
$ make build
$ make run # or ./ayame
```
2. connect to server to get web page on your browser
```
# use xdg-open in Linux instead of open of MacOS
$ [open|xdg-open] http://localhost:3000/static 
```

3. type make to do other handlings such as docker, docker-compose
```
$ make 
usage: make [build|run|kill|docker|compose|ngrok|git]
```

### API endpoints
- `/static` : file server for static web pages
    - `/upload` : directory for file upload
    - `/record` : directory for media transcoding
    - `/util` : test pages
    - `/sample` : sample media test with subtitle
- `/admin` : status monitoring
- `/event` : handle push/pull event with sub points of `/event/send|/event/recv`
- `/fetch` : fetch upload test
- `/signal` : webrtc signaling (websocket)
- `/chat` : simple chatting in a room, i.e without room concept (websocket)


### History
- 2019/12/03 : 19.05.03
    - add page and js for momo
- 2019/11/26 : 19.05.02
    - tested upload using `fetch()` in the web
- 2019/11/21 : 19.04.19
    - support `MediaSet` and it JSON
    - test subtitle display with video
    - source refactoring for file names
- 2019/11/20 : 19.04.18
    - move the directories of `upload/ & `record/` to below `asset/`
    - record the designated directory with UUID based filenames
    - define MediaSet for recording
- 2019/11/11 : 19.04.15
    - enhance uploadHandler to make a media set of upload file for service
    - modify `Dockerfile` to include ffmpeg for video conversion
- 2019/10/30 : 19.04.13
    - consider connection error and ping
    - add `isExistUser` in ayame [Commits](https://github.com/OpenAyame/ayame/commits/develop) of 2019/10/22
    - test for golang 1.12.12, 1.13.3
- 2019/10/29 : 19.04.12
    - add disk space monitoring
    - define `format={text|json|html}` in query
    - add an event channel in hub to send event to clients via sse
- 2019/10/26 : 19.04.10
    - enhance some test pages
    - set the default the number of sessions to 100
- 2019/10/25 : 19.04.09
    - define admin APIs and write its page
- 2019/10/24 : 19.04.08
    - remove `/ws` endpoint for compatibility backup
    - add `/chat` endpoint for chatting function
    - fix `docker-compose.yml` error of not assigned ports
- 2019/10/23 : 19.04.05
    - public test on admazon aws with 19.04.04
    - change `Dockerfile` to use multi-stage build to reduce docker image size
    - divide `Message` into `SignalMessage` and `ChatMessage` for its purpose 
- 2019/10/22 : 19.04.04
    - add `Dockerfile` to build its docker image, and upload it to [dockerhub.com](https://cloud.docker.com/u/agilertc/repository/docker/agilertc/ayame)
    - add `docker-compose.yml` to run with redis in a group
- 2019/10/21
    - change `http.HandleFunc` into `http.Handle` for http file server using directory mapping
    - add `sse.html` to test SSE support of server
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
- [Docker and Go Modules](https://dev.to/plutov/docker-and-go-modules-3kkn)


### Reference
- [How to format current time using a yyyyMMddHHmmss format?](https://stackoverflow.com/questions/20234104/how-to-format-current-time-using-a-yyyymmddhhmmss-format)
- [Uploading Files in Go - Tutorial](https://tutorialedge.net/golang/go-file-upload-tutorial/)
- [Basic Redis Examples with Go](https://medium.com/@gilcrest_65433/basic-redis-examples-with-go-a3348a12878e)
    - [go-redis](https://github.com/go-redis/redis):7.1k vs [redigo](https://github.com/gomodule/redigo):6.6k
    - [gilcrest/redigo-example](https://github.com/gilcrest/redigo-example)
- [Uploading files using 'fetch' and 'FormData'](https://muffinman.io/uploading-files-using-fetch-multipart-form-data/)
- [비디오에 캡션 달기 예제](http://visualssing.dothome.co.kr/temp/videocaption.html)


### Tools
- Chrome [Extenstions](https://chrome.google.com/webstore/category/extensions): [Simple WebSocket Client](https://chrome.google.com/webstore/detail/simple-websocket-client/pfdhoblngboilpfeibdedpjgfnlcodoo)
- [jrottenberg/ffmpeg](https://github.com/jrottenberg/ffmpeg) - Docker build for FFmpeg on Ubuntu / Alpine / Centos 7 / Scratch
