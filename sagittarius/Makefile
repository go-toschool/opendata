# Aiolos Makefile
# Saggitarius pod
PROGRAM="Aiolos Service"

SVC=aiolos
VERSION=v0.0.1
PROTO_DIR=proto
PROTO_SRC=$(PROTO_DIR)/*.proto
GO_PTYPES=$(SVC)/*.pb.go

BIN=$(PWD)/bin/$(SVC)
REGISTRY_URL=gcr.io/finciero-gold

clean c:
	@echo "[clean] Cleaning files..."
	@rm -f $(GO_PTYPES)
	@rm -f $(BIN)

proto p: clean
	@echo "[proto] Generating golang proto types..."
	protoc --proto_path=$(GOPATH)/src:$(PROTO_DIR)/ --go_out=plugins=grpc:$(GOPATH)/src $(PROTO_SRC)

build:
	@echo "[build] aiolos service..."
	@cd cmd/$(SVC) && GOOS=linux go build -o $(BIN) -ldflags=$(LDFLAGS)

docker: proto build
	@echo "[docker] Building image..."
	@docker build -t $(REGISTRY_URL)/$(SVC):$(VERSION) .