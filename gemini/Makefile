# Gemini Makefile
# Gemini pod
PROGRAM="Gemini Service"

VERSION=v0.0.1

SVC=gemini
SVC2=kanon
SVC3=saga

PROTO_DIR=proto
REGISTRY_URL=gcr.io/finciero-gold

PROTO_SRC_SVC=$(PROTO_DIR)/$(SVC).proto
PROTO_SRC_SVC2=$(PROTO_DIR)/$(SVC2).proto
PROTO_SRC_SVC3=$(PROTO_DIR)/$(SVC3).proto

GO_PTYPES=$(SVC)/*.pb.go
GO_PTYPES2=$(SVC2)/*.pb.go
GO_PTYPES3=$(SVC3)/*.pb.go

LDFLAGS='-extldflags "static" -X main.svcVersion=$(VERSION) -X main.svcName=$(SVC)'
LDFLAGS2='-extldflags "static" -X main.svcVersion=$(VERSION) -X main.svcName=$(SVC2)'
LDFLAGS3='-extldflags "static" -X main.svcVersion=$(VERSION) -X main.svcName=$(SVC3)'

BIN=$(PWD)/bin/$(SVC)
BIN2=$(PWD)/bin/$(SVC2)
BIN3=$(PWD)/bin/$(SVC3)

clean c:
	@echo "[clean] Cleaning files..."
	@rm -f $(GO_PTYPES)
	@rm -f $(GO_PTYPES2)
	@rm -f $(GO_PTYPES3)
	@rm -f $(BIN)

protogemini:
	@echo "[proto] Generating golang proto types..."
	protoc --proto_path=$(GOPATH)/src:$(PROTO_DIR)/ --go_out=plugins=grpc:$(GOPATH)/src $(PROTO_SRC_SVC)

protokanon:
	@echo "[proto] Generating golang proto types..."
	protoc --proto_path=$(GOPATH)/src:$(PROTO_DIR)/ --go_out=plugins=grpc:$(GOPATH)/src $(PROTO_SRC_SVC2)

protosaga:
	@echo "[proto] Generating golang proto types..."
	protoc --proto_path=$(GOPATH)/src:$(PROTO_DIR)/ --go_out=plugins=grpc:$(GOPATH)/src $(PROTO_SRC_SVC3)

proto p: clean protogemini protokanon protosaga
	@echo "[proto] Generating golang proto types..."

buildgemini:
	@echo "[build] gemini service..."
	@cd cmd/$(SVC) && GOOS=linux go build -o $(BIN) -ldflags=$(LDFLAGS)

buildkanon:
	@echo "[build] kanon service..."
	@cd cmd/$(SVC2) && GOOS=linux go build -o $(BIN2) -ldflags=$(LDFLAGS)

buildsaga:
	@echo "[build] saga service..."
	@cd cmd/$(SVC3) && GOOS=linux go build -o $(BIN3) -ldflags=$(LDFLAGS)

dockergemini:
	@echo "[docker] Building image..."
	@docker build -t $(REGISTRY_URL)/$(SVC):$(VERSION) -f DockerfileGemini .

dockerkanon:
	@echo "[docker] Building image..."
	@docker build -t $(REGISTRY_URL)/$(SVC2):$(VERSION) -f DockerfileKanon .

dockersaga:
	@echo "[docker] Building image..."
	@docker build -t $(REGISTRY_URL)/$(SVC3):$(VERSION) -f DockerfileSaga .

docker d: build dockergemini dockerkanon dockersaga
	@echo "[docker] Building image..."


build b: clean proto buildgemini buildkanon buildsaga