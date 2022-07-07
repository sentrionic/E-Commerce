.PHONY: vendor ent
vendor:
	rmdir /Q /S vendor && go get -u && go mod tidy && go mod vendor

mock:
	mockery --all --inpackage --dir=publishers

lint:
	golangci-lint run

test:
	go test -v -cover ./handler/...

ent:
	go generate ./ent