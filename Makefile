.PHONY: docs

docs: install-tools
	gomarkdoc --output ./docs/doc.md .

check-docs: install-tools
	gomarkdoc --check --output ./docs/doc.md .

install-tools:
	go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest