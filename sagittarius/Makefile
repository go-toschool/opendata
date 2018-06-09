# Aiolos Makefile
# Saggitarius pod
PROGRAM="Aiolos Service"

SVC=aiolos
PROTO_DIR=proto
PROTO_SRC=$(PROTO_DIR)/*.proto
GO_PTYPES=$(SVC)/*.pb.go

clean c:
	@echo "[clean] Cleaning files..."
	@rm -f $(GO_PTYPES)
	@rm -f $(BIN)

proto p: clean
	@echo "[proto] Generating golang proto types..."
	protoc --proto_path=$(GOPATH)/src:$(PROTO_DIR)/ --go_out=plugins=grpc:$(GOPATH)/src $(PROTO_SRC)