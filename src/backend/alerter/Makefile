.PHONY: build
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go && \
	zip main.zip main && \
	rm main

clean:
	rm main.zip