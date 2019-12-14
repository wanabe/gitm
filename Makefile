SRC=cmd/gitm/main.go cmd/log_server/main.go api/gitm/gitm.pb.go
CMDS=gitm log_server

all: $(CMDS)

.PHONY: refresh

refresh: api/gitm/gitm.pb.go ruby/lib/gitm/protobuf/gitm_pb.rb

api/gitm/%.pb.go: api/%.proto
	protoc --go_out=plugins=grpc:${GOPATH}/src $<

gitm: cmd/gitm $(SRC)
	go build ./$<

log_server: cmd/log_server $(SRC)
	go build ./$<

ruby/vendor:
	cd ruby && bundle install --path=vendor/

ruby/lib/gitm/protobuf/%_pb.rb: api/%.proto ruby/vendor
	cd ruby && bundle exec grpc_tools_ruby_protoc -I../api --ruby_out=lib/ --grpc_out=lib/ ../$<
