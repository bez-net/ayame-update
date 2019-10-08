## Ayame 기술개요

이 문서는 Ayame 신호 서버 및 데모 응용 프로그램이 어떻게 작동하는지 설명하는 것을 목적으로합니다.

### 호환성

이 신호 서버는 [webrtc/apprtc: The video chat demo app based on WebRTC](https://github.com/webrtc/apprtc)와 호환성이 있습니다.

### 신호에 대해

Ayame Server의 ws://localhost:3000/ws클라이언트에서 WebSocket 연결을 설정하고 관리하는 엔드 포인트입니다.

이 엔드 포인트에 WebSocket 연결하면 Ayame Server는 연결된 클라이언트를 보유하고 있습니다.

이 신호 서버 1:1 전용이므로 3개 이상의 클라이언트의 연결 요청을 거부합니다.

Ayame는 WebSocket 연결하는 클라이언트 중 어느 하나에서 데이터가 오면 보내는 클라이언트 이외의 연결된 모든 클라이언트에 데이터를 그대로 WebSocket으로 보냅니다. 이들은 모두 비동기식으로 이루어집니다.

이것이 "signaling"입니다.


### 연결 확립까지의 시퀀스 다이어그램

Ayame Server가 서로의 SDP 교환이나 peer connection 연결 신호에 의해 상호 작용합니다.

SDP는 WebRTC 연결에 필요한 peer connection의 내부 정보입니다.

- [RFC 4566 \- SDP: Session Description Protocol](https://tools.ietf.org/html/rfc4566)
- [Annotated Example SDP for WebRTC](https://tools.ietf.org/html/draft-ietf-rtcweb-sdp-11)

 ```

  +-------------+     +-------------------+    +-------------+
  |   browser1  |     |   Ayame Server    |    |   browser2  |
  +-----+-------+     +--------+----------+    +------+------+
        |                      |                      |
    ----------------------WebSocket 접속확립----------------
        +--------------------->|                      |
        |      websocket 접속   | <--------------------+
        |                      |     websocket 접속    |
    -----------------Peer-Connection 초기화 ---------------
        |                      |                      |
        | getUserMedia()       |                      | getUserMedia() 
        | localStream를 취득     |                      | localStream를 취득
        | peer = new PeerConnection()                 | peer = new PeerConnection()
    ----------------- 클라이언트 정보 등록 ----------------
        +--------------------->|                      | room id와 client id를 등록합니다.
        |   ws message         |                      | 2인 이하의 입실이 가능한 ayame는 accept를 반환
        |   {type: register,   |                      | 그외의 입실은 reject를 반환
        |   room_id: room_id,  |                      | TURN등 메타데이터도 여기서 교환한다.
        |   client: client_id} |                      |
        |<---------------------+                      |
        |  {type: accept }     |<---------------------|
        |                      |   ws message         | 
        |                      |    register          |  
        |                      |--------------------->|
    -----------------------SDP 교환 -----------------------
        |                      |                      |
        + peer.createOffer(),  |                      |
        | peer.setLocalDescription()                  |
        |  offerSDP를 생성       |                      |
        |                      |                      |
        +--------------------->|                      |
        |      ws message      |--------------------> |
        |      offerSDP        |   ws message         | offerSDP을 바탕으로 Remote Description을 설정
        |                      |    offerSDP          |  answerSDP을 생성하고、이를 바탕으로 localDescription을 생성한다.
        |                      |                      |　peer.setRemoteDescription(offerSDP),
        |                      |                      |  peer.createAnswer(),
        |                      | <--------------------+  peer.setLocalDescription(answer)
        | <--------------------+    ws message        |
        |      ws message      |    answerSDP         |
        |     　answerSDP      |                      |
        |                      |                      |
        + setRemoteDescription(answerSDP)             |
        | Remote Description 전송                      |
        |                      |                      |
　      |                      |                      |
    ------------------ ICE candidate를 교환 -----------------------
        |                      |                      |
　　　　 + onicecandidate() 동작 |                      |
        |  candidate 취득       |                      |
        +--------------------->+                      |
        |      ws message      +--------------------> | peerConnection에 ice candidate를 추가합니다.
        |  {type: "candidate", |   ws message         | peer.addIceCandidate(candidate)
        |    ice: candidate}   |   {type: "candidate",|
        |                      |   ice: candidate}    |　
      ==== 같은 형태로 browser2에서 browser1으로 ICE candidate 교환을 수행 ====
        |                      |                      |
     ========= ICE negotiation이 있으면 갱신 SDP를 다시 교환 =================
        |                      |                      |
        + onaddstream() 동작    |                      + onaddstream() 동작
        | remoteStream를 전송（browser2）                | remoteStream をセット(browser1)
    ------------------ Peer Connection 확립 ------------------------
 　　    |                      |                      |　
```


### 프로토콜

WS의 메시지는 JSON 형식으로 전달합니다. 모든 메시지는 속성 type을 가지고 있습니다. type다음의 5 가지입니다

- register
- accept
- reject
- offer
- answer
- candidate
- close

#### type: register

클라이언트가 Ayame Server에 room id, client id를 등록하는 메시지입니다.

```
{type: "register", "room_id" "<string>", "client_id": "<string>"}
```

이것을 받은 Ayame Server는 클라이언트가 지정한 room에 입실할 수 있는지 점검하고 가능하면 accept, 불가하다면 reject를 반환합니다.

#### type: accept

Ayame Server가 register에 담겨있는 정보를 검사하여 입실 가능하다는 것을 클라이언트에게 알리는 메시지입니다.

```
{type: "accept"}
```

향후 TURN 등의 메타 데이터도 여기에 반환 될 예정입니다. 이것을 받으면 클라이언트는 offer 교환을 시작합니다.

#### type: reject

Ayame Server가 register에 담겨있는 정보를 검사하여 입실 불가능하다는 것을 클라이언트에게 알리는 메시지입니다.

```
{type: "reject"}
```

이것을 받으면 클라이언트는 peerConnection, websocket을 닫고 초기화합니다.

#### type: offer

offer SDP를 보내는 메시지입니다.

```
{type: "offer", sdp: "v=0\r\no=- 4765067307885144980..."}
```

이를받은 클라이언트는이 SDP를 바탕으로 peer connection에 remote description을 설정합니다. 
또한 이시기에 local description을 생성하고 anwser SDP를 보냅니다.

#### type: answer

answer SDP를 보내는 메시지입니다.

```
{type: "answer", sdp: "v=0\r\no=- 4765067307885144980..."}
```

이를받은 클라이언트는이를 바탕으로 peer connection에 remote description을 설정합니다.

#### type: candidate

ice candidate를 교환하는 메시지입니다.

```
{type: "candidate", ice: {candidate: "...."}}
```

이를받은 클라이언트는 peer connection에 ice candidate를 추가합니다.

#### type: close

peer connection을 끊어 졌음을 알리는 메시지입니다.

```
{type: "close"}
```

이를받은 클라이언트는 peer connection을 닫고 원격 (수신자)의 video element를 파기합니다.

