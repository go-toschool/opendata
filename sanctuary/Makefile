# Sanctuary Makefile
# Sanctuary pod
PROGRAM="Sanctuary Service"

VERSION=v0.0.1

SVC=sanctuary

PROTO_DIR=proto
REGISTRY_URL=gcr.io/finciero-gold

PROTO_SRC_SVC=$(PROTO_DIR)/$(SVC).proto

GO_PTYPES=$(SVC)/*.pb.go

LDFLAGS='-extldflags "static" -X main.svcVersion=$(VERSION) -X main.svcName=$(SVC)'

BIN=$(PWD)/bin/$(SVC)

clean c:
	@echo "[clean] Cleaning files..."
	@rm -f $(BIN)

build: clean
	@echo "[build] service..."
	@cd cmd/$(SVC) && GOOS=linux go build -o $(BIN) -ldflags=$(LDFLAGS)

docker: build
	@echo "[docker] Building image..."
	@docker build -t $(REGISTRY_URL)/$(SVC):$(VERSION) -f manifests/$(SVC) .

