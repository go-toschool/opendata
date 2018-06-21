# Capricornius Makefile
# Capricornius pod
PROGRAM="Capricornius Service"

VERSION=v0.0.1

SVC=shura

PROTO_DIR=proto
REGISTRY_URL=gcr.io/finciero-gold

PROTO_SRC_SVC=$(PROTO_DIR)/$(SVC).proto

GO_PTYPES=$(SVC)/*.pb.go

LDFLAGS='-extldflags "static" -X main.svcVersion=$(VERSION) -X main.svcName=$(SVC)'

BIN=$(PWD)/bin/$(SVC)

clean c:
	@echo "[clean] Cleaning files..."
	@rm -f $(GO_PTYPES)
	@rm -f $(BIN)

proto: clean
	@echo "[proto] Generating golang proto types..."
	protoc --proto_path=$(GOPATH)/src:$(PROTO_DIR)/ --go_out=plugins=grpc:$(GOPATH)/src $(PROTO_SRC_SVC)

build: proto
	@echo "[build] gemini service..."
	@cd cmd/$(SVC) && GOOS=linux go build -o $(BIN) -ldflags=$(LDFLAGS)

docker: build
	@echo "[docker] Building image..."
	@docker build -t $(REGISTRY_URL)/$(SVC):$(VERSION) -f manifests/$(SVC) .

