
-include local.Makefile

.PHONY: lint
lint: ## lint
	golangci-lint run

.PHONY: fmt
fmt: tidy ## tidy,format and imports
	gofumpt -w `find . -type f -name '*.go' -not -path "./vendor/*"`
	goimports -w `find . -type f -name '*.go' -not -path "./vendor/*"`

.PHONY: tidy
tidy: ## go mod tidy
	go mod tidy

install:
	CGO_ENABLED=0 go install ./cmd/torrenti

.PHONY: bin/torrenti
bin/torrenti:
	GOAMD64=v4 CGO_ENABLED=0 go build -ldflags '-s -w' -trimpath -o bin/torrenti ./cmd/torrenti

.PHONY: bin/torrenti_linux_amd64
bin/torrenti_linux_amd64:
	GOOS=linux GOAMD64=v3 CGO_ENABLED=0 go build -ldflags '-s -w' -trimpath -o bin/torrenti_linux_amd64 ./cmd/torrenti

#ls -d pkg/indexer/plugins/* | xargs -n 1 -I {} sh -c 'CGO_ENABLED=1 go build -ldflags "-s -w" -trimpath -buildmode=plugin -o bin/plugins/`basename {}`.so ./{}'
build:
	go env
	CGO_ENABLED=0 go build -ldflags '-s -w' -trimpath -o bin/torrenti ./cmd/torrenti
