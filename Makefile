SHELL=/bin/bash
GITSHORTHASH=$(shell git log -1 --pretty=format:%h)
GO ?= go
VERSION ?= $(GITSHORTHASH)

# Insert build metadata into binary
LDFLAGS := -X github.com/attadanta/tado-metrics/cmd.TadoMetricsVersion=$(VERSION)
LDFLAGS += -X github.com/attadanta/tado-metrics/cmd.TadoMetricsGitCommit=$(GITSHORTHASH)

.PHONY: build
build:
	env GOOS=linux GOARCH=arm GOARM=5 $(GO) build -ldflags "$(LDFLAGS)" -o "tado-metrics" main/tado.go

.PHONY: package
package: build
	tar -cvzf tado-metrics.tar.gz tado-metrics service.env INSTALL.sh

.PHONY: clean
clean:
	rm -rf tado-metrics.tar.gz tado-metrics
