## Ayame를 사용해보자.

이리 엔트리는 크롤링합니다. 디렉토리 구성은 다음과 같게됩니다.

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

`make` 로드 빌드가 됩니다.

```
$ make
```

## Server 기동

빌라를 성공한 후에, 다음의 명령으로 아야메 서버를 작동시킬 수 있습니다.

```
$ ./ayame
```
기동한 후에, http://localhost:3000 로 액세스할 때 데모 화면에 액세스 할 수 있습니다.

アクセス時に各ブラウザで「カメラ・マイクでのアクセス」権限を要求された場合は「許可する」を選択してください。

権限を確認できたら、「接続する」を選択してください。

別のタブ or ブラウザから同様にアクセスして、互いの画面が表示されたら接続成功です。

※ あくまで Peer 2 Peer なので、最大 2 クライアントまでの接続しかできません。

切断するときは「切断する」を選択してください。


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

## `over_ws_ping_pong` オプションについて

- `config.yaml` にて `over_ws_ping_pong: true` に設定した場合、 ayame はクライアントに対して(WebSocket の ping frame の代わりに) ** 9 ** 秒おきに JSON 形式で `{"type": "ping"}` メッセージを送信します。
- これに対してクライアントは ** 10 ** 秒以内に JSON 形式で `{"type": "pong"}` を返すことで ping-pong を実現します。

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


## `register` メッセージについて

クライアントは ayame への接続可否を問い合わせるために WebSocket に接続した際に、まず `"type": "register"` のJSON メッセージを WS で送信する必要があります。
register で送信できるプロパティは以下になります。


- `"type"`: (string): 必須。 `"register"` を指定する。
- `"clientId"`: (string): 必須
- `"roomId": (string): 必須
- `"key"`(string): Optional
- `"authnMetadata"`(Object): Optional
    - 多段ウェブフック認証の際に利用することができます。多段ウェブフック認証については後述します。


## `auth_webhook_url` オプションについて

`config.yaml` にて `auth_webhook_url` を指定している場合、 ayame は client が {"type": "register" } メッセージを送信してきた際に `config.yaml` に指定した `auth_webhook_url` に対して認証リクエストをJSON 形式で POST します。


このとき、{"type": "register" } のメッセージに

- `"key"`(string)
- `"room_id"`: (string)

を含めていると、そのデータを ayame はそのまま指定した `auth_webhook_url` に JSON 形式で送信します。


また、 auth webhook の返り値は JSON 形式で、以下のように想定されています。

- `allowed`: boolean。認証の可否 (必須)
- `reason`: string。認証不可の際の理由 (`allowed` が false の場合のみ必須)
- `auth_webhook_url`: 多段認証用の webhook url。(optional、多段認証をしない場合不要)
    - 多段認証については次の項で説明します。

`allowed` が false の場合 client の ayame への WebSocket 接続は切断されます。

この auth_webhook はシグナリング key とroom ID の結びつきを確認する想定のものです。


### 多段ウェブフック認証について

`auth_webhook_url` を指定して、その `auth_webhook_url` からの返り値の JSON プロパティに `auth_webhook_url` が指定してある場合、
ayame は通常の認証 wehbook での認証後:wその URL に対してさらに認証リクエストを POST します。
この `auth_webhook_url` へのリクエスト、レスポンスは以下のように想定されています。

#### request

- `host`: string。クライアントの host。
- `authn_metadata`(Object)
    - register 時に `authn_metadata` をプロパティとして指定していると、その値がそのまま付与されます。


#### response

- `allowed`: boolean。認証の可否 (必須)
- `reason`: string。認証不可の際の理由 (`allowed` が false の場合のみ)
- `authzMetadata`(Object, Optional)
    - 任意に払い出せるメタデータ。 client はこの値を読み込むことで、例えば username を認証サーバから送ったりということも可能になる。


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

