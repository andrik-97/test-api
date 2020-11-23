GOFMT_FILES = $(shell go list -f '{{.Dir}}' ./... | grep -v '/pb' | grep -v '/postgres/model')

.PHONY: fmtcheck
fmtcheck:
	@command -v goimports > /dev/null 2>&1 || go get golang.org/x/tools/cmd/goimports
	@CHANGES="$$(goimports -d $(GOFMT_FILES))"; \
		if [ -n "$${CHANGES}" ]; then \
			echo "Unformatted (run goimports -w .):\n\n$${CHANGES}\n\n"; \
			exit 1; \
		fi
	@# Annoyingly, goimports does not support the simplify flag.
	@CHANGES="$$(gofmt -s -d $(GOFMT_FILES))"; \
		if [ -n "$${CHANGES}" ]; then \
			echo "Unformatted (run gofmt -s -w .):\n\n$${CHANGES}\n\n"; \
			exit 1; \
		fi

.PHONY: spellcheck
spellcheck:
	@command -v misspell > /dev/null 2>&1 || go get github.com/client9/misspell/cmd/misspell
	@misspell -locale="US" -error -source="text" **/*

.PHONY: staticcheck
staticcheck:
	@command -v staticcheck > /dev/null 2>&1 || go get honnef.co/go/tools/cmd/staticcheck
	@staticcheck -checks="all" -tests $(GOFMT_FILES)

.PHONY: test
test:
	@go test \
		-cover \
		$(GOFMT_FILES)
