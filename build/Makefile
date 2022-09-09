VERSION=1.0.1
DATETIME=`date +%FT%T%z`

GOFILES=`find . -name "*.go" -type f`
MK := ./make.sh
list:
	echo ${VERSION}
	echo ${GOFILES}

# gofmt
# -s	simplify code
# -w	write result to (source) file instead of stdout
# -d	display diffs instead of rewriting files
fmt:
	gofmt -s -w ${GOFILES}

fmt-check:
	diff=?(gofmt -s -d $(GOFILES)); \
	if [-n "$$diff"];then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

test:
	go test
build: 
	${MK} goweb gw
# include ./*.MK
# build not represend physical command
# .PHONY: build
