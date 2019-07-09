IMAGE := bluex/pgexport
CMDS := $(shell find cmd -mindepth 1 -maxdepth 1 -type d | sed 's/ /\\ /g' | xargs -n1 basename)
BINS := $(patsubst %,bin/%,$(CMDS))
GOFILES := $(shell find . -type f -name '*.go' -not -path './vendor/*' -not -path './cmd/*')
$(foreach cmd,$(CMDS),\
	$(eval GOFILES_$(cmd) := $(shell find ./cmd/$(cmd) -type f -name '*.go')))
PROTO_DIR := protobuf
PROTO_FILES := $(wildcard $(PROTO_DIR)/**/*.proto)
PROTO_GO_FILES := $(PROTO_FILES:.proto=.pb.go)
PROTO_GO_DIR := pkg/proto
PROTO_GO_FILES := $(patsubst $(PROTO_DIR)/%.proto,$(PROTO_GO_DIR)/%.pb.go,$(PROTO_FILES))
PLATFORMS := darwin linux windows
ARCHITECTURES := 386 amd64
CROSS := $(foreach p,$(PLATFORMS),$(patsubst %,$p-%,$(ARCHITECTURES)))
$(info "cross: $(CROSS)")
#$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES),\
	$(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -o ./bin/$(@F)-$(GOOS)-$(GOARCH) ./cmd/$(@F))))

.PHONY: $(CMDS) clean distclean

.DEFAULT: default

.PHONY: default
default: $(CMDS)

$(CMDS): $(BINS)

.PHONY: pb-sync
pb-sync:
	@$(MAKE) -C $(PROTO_DIR) sync

.PHONY: pb
pb: $(PROTO_GO_FILES)

.PHONY: mock
mock: $(MOCK_FILES)

$(PROTO_GO_DIR)/%.pb.go: $(PROTO_DIR)/%.proto
	@mkdir -p $(PROTO_GO_DIR)
	protoc -I $(PROTO_DIR) --go_out=paths=source_relative,plugins=grpc:$(PROTO_GO_DIR) $<

$(MOCK_FILES): test/mock/% : pkg/%
	@mkdir -p test/mock
	mockgen -source $< -destination $@

.SECONDEXPANSION:
bin/%: export GOOS := $(GOOS)
bin/%: export GOARCH := $(GOARCH)
bin/%: $(PROTO_GO_FILES) $(GOFILES) $$(GOFILES_$$(@F))
	go build -o ./bin/$(@F) ./cmd/$(@F)

.PHONY: vet
vet:
	@go vet ./...

.PHONY: lint
lint:
	@golint $(shell go list ./...)

.PHONY: fmt
fmt:
	@go fmt ./...
	@-goimports -w $$(go list -f {{.Dir}} ./... | grep -v "/pkg\/proto/\|/test\/mock/")

.PHONY: pb-fmt
pb-fmt:
	@-clang-format -style='{ColumnLimit: 0}' -i $(PROTO_FILES)

.PHONY: test
test: pb mock
	@go clean -testcache
	@go test -p 1 -v -cover ./...

.PHONY: docker
docker:
	@docker build . -t $(IMAGE)

.PHONY: clean
clean:
	@rm -rf bin

.PHONY: distclean
distclean: clean
	@rm -f $(PROTO_GO_FILES)
	@rm -f $(MOCK_FILES)

