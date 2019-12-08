SRC=cmd/gitm/main.go api/gitm/gitm.pb.go
EXE=gitm

all: $(EXE)

api/gitm/%.pb.go: api/%.proto
	protoc --go_out=${GOPATH}/src $<

gitm: cmd/gitm $(SRC)
	go build ./$<
