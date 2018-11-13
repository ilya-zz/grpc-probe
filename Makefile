APIDIR = api
PYDIR = python
PROTOFILE := $(APIDIR)/api.proto
PROTO_INC := ${GOPATH}/src/github.com/gogo/protobuf/protobuf

PB := $(api)/api.pb.go

BINDIR=build
SERVER=$(BINDIR)/server
CLIENT=$(BINDIR)/client

SRV_PORT := $(or ${SP}, 7777)

.PHONY: gogo protoclean proto bench all run dirs

dirs:
	@mkdir -p $(BINDIR)
	@mkdir -p $(APIDIR)

gogo: dirs protoclean 	
	protoc -I=. -I=$(PROTO_INC) --gogoslick_out=plugins=grpc:. --python_out=. $(PROTOFILE)
	python -m grpc_tools.protoc -Iapi --python_out=$(PYDIR) --grpc_python_out=$(PYDIR) $(PROTOFILE)

proto: dirs protoclean	
	protoc -I=. --go_out=plugins=grpc:. $(PROTOFILE)
	python -m grpc_tools.protoc -Iapi --python_out=$(PYDIR) --grpc_python_out=$(PYDIR) $(PROTOFILE)

bench%:
	go test -bench=. -benchmem -v -benchtime 3s .

protoclean:
	@rm -f $(APIDIR)/*.pb.go

$(PB): gogo

$(SERVER): $(PB) cmd/server/main.go
	go build -o $@ ./cmd/server

$(CLIENT): $(PB) cmd/client/main.go
	go build -o $@ ./cmd/client	

build: $(SERVER) $(CLIENT)

clean:
	@rm -f $(CLIENT) $(SERVER)	

stops:
	@killall -9 server || true

runs: $(SERVER) stops
	./$(SERVER) -p $(SRV_PORT) &

run: build runs	
	@./$(CLIENT)

pyth: runs
	python $(PYDIR)/pclient.py