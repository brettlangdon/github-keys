github-keys: github-keys.go
	go build

run:
	go run ./github-keys.go

clean:
	rm -f ./github-keys

.PHONY: run clean
