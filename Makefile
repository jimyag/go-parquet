PACKAGES=`go list ./... | grep -v example`

test:
	GORACE=atexit_sleep_ms=0 go test -trimpath -failfast -race -cover  ${PACKAGES}

fmt:
	go fmt ./...

PACKAGES=`go list ./... | grep -v -e example -e "parquet/parquet"`
staticcheck:
	staticcheck ${PACKAGES}

.PHONEY: test
