# WebRTC Signaling Server Ayame

[! [GitHub tag (latest SemVer) (https://img.shields.io/github/tag/OpenAyame/ayame.svg) (https://github.com/OpenAyame/ayame)
[! [License] (https://img.shields.io/badge/License-Apache%202.0-blue.svg) (https://opensource.org/licenses/Apache-2.0)
[! [Actions Status (https://github.com/OpenAyame/ayame/workflows/Go%20Build%20&%20Format/badge.svg) (https://github.com/OpenAyame/ayame/actions)

## WebRTC Signaling Server Ayame 대해

WebRTC Signaling Server Ayame는 WebRTC를위한 시그널링 서버입니다.

WebRTC의 P2P에서만 작동합니다. 또한 동작을 1 룸을 최대 2 명으로 제한하여 코드를 작게 유지하고 있습니다.

AppRTC 호환 룸 기능을 가지고 있으며, 객실 수는 서버 사양에 따라 다르지만 1 만까지 처리 할 수 ​​있도록되어 있습니다.

## OpenAyame 프로젝트에 대해

OpenAyame 프로젝트는 WebRTC Signaling Server Ayame를 오픈 소스로 공개하고 지속적으로 개발함으로써 WebRTC를 배우기 쉽고 프로젝트입니다.

자세한 내용은 아래를 참조하십시오.

[OpenAyame 프로젝트 (http://bit.ly/OpenAyame)

## 개발에 대해

Ayame는 오픈 소스 소프트웨어이지만 개발에 오픈하지 않습니다.
따라서 의견과 풀 요청을 받아도 즉시 채택하지 않습니다.

우선 Discord로 연락하십시오.

##주의

- Ayame는 P2P 밖에 대응하지 않습니다
- Ayame 1 실 최대 2 명까지 밖에 대응하고 있지 않습니다
- 샘플이 이용하고있는 STUN 서버는 Google의 것을 이용하고 있습니다

## 사용해보기

Ayame를 사용해보고 싶은 사람은 [USE.md (doc / USE.md)를 읽어 보시기 바랍니다.

## SDK를 사용해 본다

쉽게 Ayame를 이용할 수있는 Web SDK를 제공하고 있습니다.

[OpenAyame / ayame \ -web \ -sdk : Web SDK for WebRTC Signaling Server Ayame (https://github.com/OpenAyame/ayame-web-sdk)

```javascript
const conn = Ayame.connection ( 'wss : //example.com : 3000 / signaling', 'test-room');
const startConn = async () => {
    const mediaStream = await navigator.mediaDevices.getUserMedia ({audio : true, video : true});
    await conn.connect (mediaStream);
    conn.on ( 'disconnect'(e) => console.log (e));
    conn.on ( 'addstream'(e) => {
        document.querySelector ( '# remote-video'). srcObject = e.stream;
    });
    document.querySelector ( '# local-video'). srcObject = mediaStream;
};
startConn ();
```

## React 샘플을 사용해 보는

**이 저장소에있는 샘플과 똑같은 동작이되어 있습니다 **

[OpenAyame / ayame \ -react \ -sample (https://github.com/OpenAyame/ayame-react-sample)

## WebRTC 신호 서비스 Ayame Lite를 사용해 본다

Ayame을 이용한 신호 서비스를 제공하고 있습니다.

```
wss : //ayame-lite.shiguredo.jp/signaling
```

인증 등은 현재 걸고 있지 않으므로, 룸 ID는 다른 사람으로부터 추측 할 수없는 값을 사용하도록하십시오.

자세한 내용은 다음을 참조하십시오.

[WebRTC 신호 서비스 Ayame Lite 개발 로그 (https://gist.github.com/voluntas/396167bd197ba005ae5a9e8c5e60f7cd)

## 구조의 자세한 정보를 원하시면

Ayame 자세한 내용을 알고 싶은 사람은 [DETAIL.md (doc / DETAIL.md)를 읽어 보시기 바랍니다.

## 관련 제품

[hakobera / serverless-webrtc-signaling-server (https://github.com/hakobera/serverless-webrtc-signaling-server)가 Ayame 호환 서버로 공개 / 개발되고 있습니다. AWS에 의해 서버 주소를 실현 한 WebRTC P2P Signaling Server입니다.

## 라이센스

Apache License 2.0

```
Copyright 2019, Shiguredo Inc, kdxu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS"BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## 지원

WebRTC Signaling Server Ayame 버그보고는 GitHub Issues로 부탁드립니다. 그렇지 대해서는 Discord에 부탁합니다.

### 버그보고

https://github.com/OpenAyame/ayame/issues

### Discord

최선형 운영하고 있습니다.

https://discord.gg/mDesh2E

### 유료 지원

** 時雨堂에서는 유료 지원은 제공하고 있지 않습니다 **

- [kdxu \ (Kyoko KADOWAKI \) (https://github.com/kdxu)가 유료로 지원 및 사용자 정의를 제공합니다. Discord 통해 @kdxu에 연락을 부탁드립니다.