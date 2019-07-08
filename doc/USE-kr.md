## Ayame를 사용해보자.

이 레포지토리를 크롤링해보면 디렉토리 구성은 다음과 같습니다.

```
$ ./
.
├── sample/
│   ├── index.html
│   ├── main.css
│   └── webrtc.js
├── .doc/
│   ├── USE.md
│   └── DETAIL.md
├── go.mod
├── go.sum
├── ws_handler.go
├── client.go
├── hub.go
└── main.go
```


## Go Install

추천 버전은 다음과 같습니다.
```
go 1.12
```

## visualing

```
$ go build
```

`make` 로도 빌드가 됩니다.

```
$ make
```

## Server 기동

빌라를 성공한 후에, 다음의 명령으로 아야메 서버를 작동시킬 수 있습니다.

```
$ ./ayame
```
기동한 후에, http://localhost:3000 로 액세스할 때 데모 화면에 액세스할 수 있습니다.

액세스 할 때 각 브라우저에 "카메라 마이크에 액세스"권한을 요청한 경우 "허용"을 선택하십시오.

권한을 확인한 후 "연결"을 선택하십시오.

다른 탭이나 브라우저에서 동일하게 액세스하여 서로의 화면이 나타나면 연결 성공입니다.

※ 어디까지나 Peer 2 Peer 용이므로、최대 2명의 클라이언트만 접속할 수 있습니다.

연결을 끊고자 할때는 "절단한다"고 선언해주세요.


## Comment

```
$ ./ayame version
WebRTC Signaling Server Ayame version 19.02.1⏎
```

```
$ ./ayame -c ./config.yaml
time="2019-06-10T00:23:16+09:00" level=info msg="Setup log finished."
time="2019-06-10T00:23:16+09:00" level=info msg="WebRTC Signaling Server Ayame. version=19.02.1"
time="2019-06-10T00:23:16+09:00" level=info msg="running on http://localhost:3000 (Press Ctrl+C quit)"
```

```
$ ./ayame -help
Usage of ./ayame:
  -c string
    	ayame configuration file (yaml) (default "./config.yaml")
```

## `over_ws_ping_pong` 옵션에 대하여

- `config.yaml`에서`over_ws_ping_pong : true`로 설정 한 경우 ayame는 클라이언트에 대해 (WebSocket의 ping frame 대신) ** 9 ** 초 소리>와에 JSON 형식의`{ "type ":"ping "}`메시지를 보냅니다.
- 이에 대해 클라이언트 ** 10 ** 초 이내에 JSON 형식의`{ "type": "pong"}`을 반환에서 ping-pong을 제공합니다.

클라이언트 (자바 스크립트)의 샘플 코드를 아래에 설명합니다.

```javascript
ws = new WebSocket(signalingUrl);
ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      console.log(message.type)
      switch(message.type){
        case 'ping': {
          console.log('Received Ping, Send Pong.');
          ws.send(JSON.stringify({
            "type": "pong"
          }))
          break;
        }
        ...
```


## `register` 메시지에 대하여

클라이언트는 ayame 연결 여부를 문의하기 위해 WebSocket에 연결했을 때 먼저` "type": "register"`의 JSON 메시지를 WS로 전송해야합니다.
필요가 있습니다.
register로 보낼 수 속성은 다음입니다.


- `"type"`: (string): 필수. `"register"` 을 지정한다.
- `"clientId"`: (string): 필수 
- `"roomId": (string):  필수
- `"key"`(string): Optional
- `"authnMetadata"`(Object): Optional
    - 단 웹 후크 인증시 사용할 수 있습니다. 다단 웹 후크 인증에 대해서는 후술합니다.


## `auth_webhook_url` 옵션에 대하여

`config.yaml` にて `auth_webhook_url` を指定している場合、 ayame は client が {"type": "register" } メッセージを送信してきた際に `config.yaml` に指定した `auth_webhook_url` に対して認証リクエストをJSON 形式で POST します。


이 때, { "type": "register"} 메시지에

- `"key"`(string)
- `"room_id"`: (string)

を含めていると、そのデータを ayame はそのまま指定した `auth_webhook_url` に JSON 形式で送信します。


また、 auth webhook の返り値は JSON 形式で、以下のように想定されています。

- `allowed`: boolean。認証の可否 (必須)
- `reason`: string。認証不可の際の理由 (`allowed` が false の場合のみ必須)
- `auth_webhook_url`: 多段認証用の webhook url。(optional、多段認証をしない場合不要)
    - 多段認証については次の項で説明します。

`allowed` が false の場合 client の ayame への WebSocket 接続は切断されます。
 
이 auth_webhook는 신호 key와 room ID의 관계를 확인하는 것을 가정한 것입니다.


### 다단계 webhook 설정에 대하여

`auth_webhook_url` を指定して、その `auth_webhook_url` からの返り値の JSON プロパティに `auth_webhook_url` が指定してある場合、
ayame は通常の認証 wehbook での認証後:wその URL に対してさらに認証リクエストを POST します。
この `auth_webhook_url` へのリクエスト、レスポンスは以下のように想定されています。

#### request

- `host`: string. 클라이언트의 host.
- `authn_metadata`(Object)
    - register시`authn_metadata`를 속성으로 지정하고, 그 값이 그대로 부여됩니다.


#### response

- `allowed`: boolean. 허용가부 (필수)
- `reason`: string. 허용불가의 이유 (allowed가 false가 되는 이유)
- `authzMetadata`(Object, Optional)
    - 임의로 들어가고 나가는 메타 데이터. client는이 값을 가져 와서 예를 들어 username을 인증 서버에서 보내거나하는 것도 가능하게된다.

```
{"allowed": true, "authzMetadata": {"username": "kdxu", "owner": "true"}}
```

이 다단계 auth_webhook는 이용자가 지정 한 인증 웹 페이지 URL를 이용하기 위한 것으로서 상정하고 있습니다.


### local에서 wss/https를 시험해보는 경우 

[ngrok - secure introspectable tunnels to localhost](https://ngrok.com/) 사용 확실히하는 것이 좋습니다.

```
$ ngrok http 3000
ngrok by @xxxxx
Session Status online
Account        xxxxx
Forwarding     http://xxxxx.ngrok.io -> localhost:3000
Forwarding     https://xxxxx.ngrok.io -> localhost:3000
```

