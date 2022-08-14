build: 
	go build -o gop2p .

test:
	make build
	GOOS=linux GOARCH=arm64 go build -o gop2p-arm64 .
	scp gop2p* vps.simbafs.cc:~/
	rm gop2p-arm64
