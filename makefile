update-pkg-cache:
	mkdir -p /tmp/go-temp && \
	cd /tmp/go-temp && \
	go mod init temp && \
	GOPROXY=https://proxy.golang.org GO111MODULE=on \
	go get github.com/datamango-uk/montecarlo@v$(VERSION) && \
	rm -rf /tmp/go-temp