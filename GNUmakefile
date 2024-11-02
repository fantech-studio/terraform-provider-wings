default: testacc

.PHONY: testacc
testacc:
	TF_ACC=1 \
	EDGE_API_KEY_ID=$(EDGE_API_KEY_ID) \
	EDGE_API_KEY=$(EDGE_API_KEY) \
	EDGE_API_ENDPOINT=$(EDGE_API_ENDPOINT) \
	go test -race -parallel 3 ./... -v $(TESTARGS) -timeout 120m

.PHONY: install
install: VERSION ?= 0.0.1
install: PLUGIN_ARCH ?= darwin_arm64
install:
	go build -o bin/terraform-provider-wings
	mkdir -p ~/workspace/terraform-provider-wings
	cp bin/terraform-provider-wings ~/workspace/terraform-provider-wings
