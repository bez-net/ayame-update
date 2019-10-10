# 릴리스 노트

- UPDATE
    - 호환성이있는 변경
- ADD
    - 호환성이 추가
- CHANGE
    - 호환성이없는 변경
- FIX
    - 버그 수정


## develop

- [ADD] CI의 go를 1.13에 올리기
- [UPDATE] @kdxu 권장 go version을 1.13 올린다

## 19.08.0

2019-08-16

- [UPDATE]`/ ws` 엔드 포인트와 동일한 것을`/ signaling` 엔드 포인트로 추가
- [UPDATE] ayame register시 key 보낼 수 있도록
- [UPDATE] auth webhook의 결과에 iceServers이 있으면 반환하도록

## 19.07.1
- [CHANGE] @kdxu 샘플을 ayame web-sdk를 이용한 것으로 대체

## 19.07.0

- [UPDATE] @kdxu -overWsPingPong 옵션 over WS의 ping-pong에도 대응할 수 있도록했다
- FIX : @kdxu 샘플을 unified plan에 해당하는
- [ADD] @kdxu ayame 시작할 때 조금 설명을 낸다
- [ADD] @kdxu`ayame version` 버전을 표시하도록
- [ADD] @kdxu 인증 웹 연결 기능을 추가하는
- [ADD] @kdxu 다단 인증 웹 연결 기능을 추가하는
- [CHANGE] @kdxu 설정을`config.yaml`에 새기는 것처럼 변경


## 19.02.1

- FIX : @kdxu uuid를 사용하지 않고, client_id으로 잡고 돌리도록 수정

## 19.02.0

** 첫 출시 **

- [ADD] AppRTC 호환 신호 서버 추가
- [ADD] 룸 기능 추가
- [ADD] type : accept / reject 추가
- [ADD] 3 명 이상은 킥하는 기능 추가