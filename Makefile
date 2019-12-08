SRC=gitm.go pb/gitm.pb.go
EXE=gitm

all: gitm

pb/%.pb.go: %.proto
	protoc --go_out=pb $<

gitm: $(SRC)
	go build
