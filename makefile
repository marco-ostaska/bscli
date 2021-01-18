
tmp:
	mkdir /tmp/bscli

compile:
	echo "Compiling for windows"
	env GOOS=windows GOARCH=amd64 go build -o /tmp/bscli/bscli-windows-${BSCLI_VERSION}.exe -ldflags "-X github.com/marco-ostaska/bscli/cmd.version=${BSCLI_VERSION}"
	echo "Compiling for linux"
	go build -o /tmp/bscli/bscli-linux-${BSCLI_VERSION} -ldflags "-X github.com/marco-ostaska/bscli/cmd.version=${BSCLI_VERSION}"
	echo "Compiling for linux"
	echo "Compiling for osx"
	env GOOS=darwin GOARCH=amd64 go build -o /tmp/bscli/bscli-darwin-${BSCLI_VERSION} -ldflags "-X github.com/marco-ostaska/bscli/cmd.version=${BSCLI_VERSION}"
	echo "Compiling for linux"
	echo "go install local"
	go install -ldflags "-X github.com/marco-ostaska/bscli/cmd.version=${BSCLI_VERSION}"
	echo "Compiling for linux"
