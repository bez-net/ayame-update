PROG=ayame
VERSION=19.02.1

usage:
	@echo "make [build|run|kill|ngrok|git]"

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
	./$(PROG) &

kill k:
	pkill $(PROG)

log:
	tail -f $(PROG).log

# -----------------------------------------------------------------------------------------------------------------------
ngrok:
	@echo "make (ngrok) [install|run]"

ngrok-install ni:
	snap install ngrok

ngrok-run nr:
	ngrok http 3000

#-----------------------------------------------------------------------------------------
git g:
	@echo "make (git) [update|login|tag|status]"

git-update gu:
	git add .gitignore *.md Makefile doc/ sample/
	#git commit -m "initial commit"
	#git remote remove go.mod sse.go
	#git commit -m "add examples"
	git commit -m "update contents"
	git push

git-login gl:
	git config --global user.email "sikang99@gmail.com"
	git config --global user.name "Stoney Kang"
	git config --global push.default matching
	#git config --global push.default simple
	git config credential.helper store

git-tag gt:
	git tag v0.0.3
	git push --tags

git-status gs:
	git status
	git log --oneline -5
#-----------------------------------------------------------------------------------------
