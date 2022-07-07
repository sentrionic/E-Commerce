.PHONY: vendor ent
vendor:
	rmdir /Q /S vendor && go get -u && go mod tidy && go mod vendor

mock:
	mockery --all --inpackage --dir=publishers

lint:
	golangci-lint run
