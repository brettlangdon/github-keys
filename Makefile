github-keys: github-keys.go
	go build

run:
	go run ./github-keys.go -username ${USERNAME}

clean:
	rm -f ./github-keys

.PHONY: run clean
