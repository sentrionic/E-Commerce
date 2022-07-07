.PHONY: vendor
vendor:
	rmdir /Q /S vendor && go get -u && go mod tidy && go mod vendor

lint:
	golangci-lint run

test:
	go test -v -cover ./handler/...