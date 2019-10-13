#
# Makefile for ayame, WebRTC signaling server
#
PROG=ayame
VERSION=19.02.8
# -----------------------------------------------------------------------------------------------------------------------
usage:
	@echo "WebRTC signaling server : $(PROG) $(VERSION)"
	@echo "> make [build|run|kill|ngrok|git]"

# -----------------------------------------------------------------------------------------------------------------------
build b: *.go
	GO111MODULE=on go build -ldflags '-X main.AyameVersion=${VERSION}' -o $(PROG)

build-darwin bd: *.go
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -ldflags '-X main.AyameVersion=${VERSION}' -o bin/$(PROG)-darwin
build-linux bl: *.go
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags '-s -w -X main.AyameVersion=${VERSION}' -o bin/$(PROG)-linux

check:
	GO111MODULE=on go test ./...

fmt:
	go fmt ./...

clean:
	rm -rf $(PROG)

run r:
	./$(PROG)

kill k:
	pkill $(PROG)

log l:
	tail -f $(PROG).log
# -----------------------------------------------------------------------------------------------------------------------
ngrok n:
	@echo "> make (ngrok) [install|run]"

ngrok-install ni:
	snap install ngrok

ngrok-run nr:
	ngrok http 3000
#-----------------------------------------------------------------------------------------
open o:
	@echo "> make (open) [orig|page|app]"

open-orig oo:
	xdg-open https://github.com/OpenAyame/ayame

open-page op:
	xdg-open https://github.com/sikang99/ayame

open-app oa:	# AppRTC
	xdg-open https://github.com/webrtc/apprtc
#-----------------------------------------------------------------------------------------
git g:
	@echo "> make (git) [update|login|tag|status]"

git-update gu:
	git add .gitignore *.md Makefile doc/ sample/ go.* *.go *.yaml certs/
	#git remote remove go.mod sse.go
	git commit -m "data structures are modified for client, room, hub"
	git push

git-login gl:
	git config --global user.email "sikang99@gmail.com"
	git config --global user.name "Stoney Kang"
	git config --global push.default matching
	git config credential.helper store

git-tag gt:
	git tag v0.0.3
	git push --tags

git-status gs:
	git status
	git log --oneline -5
#-----------------------------------------------------------------------------------------
