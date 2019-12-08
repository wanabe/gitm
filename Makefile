SRC=cmd/gitm/main.go api/gitm/gitm.pb.go
EXE=gitm

all: $(EXE)

.PHONY: refresh

refresh: api/gitm/gitm.pb.go ruby/lib/gitm/protobuf/gitm_pb.rb

api/gitm/%.pb.go: api/%.proto
	protoc --go_out=${GOPATH}/src $<

gitm: cmd/gitm $(SRC)
	go build ./$<

ruby/vendor:
	cd ruby && bundle install --path=vendor/

ruby/lib/gitm/protobuf/%_pb.rb: api/%.proto ruby/vendor
	cd ruby && bundle exec grpc_tools_ruby_protoc -I../api --ruby_out=lib/ --grpc_out=lib/ ../$<
