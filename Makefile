PACKAGES=`go list ./... | grep -v example`

test:
	GORACE=atexit_sleep_ms=0 go test -trimpath -failfast -race -cover  ${PACKAGES}

format:
	go fmt github.com/jimyag/parquet-go/...

.PHONEY: test
