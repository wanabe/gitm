require "gitm"

stub = Gitm::Protobuf::Log::Stub.new('localhost:50051', :this_channel_is_insecure)
commit = stub.get(Gitm::Protobuf::LogRequest.new(object: Gitm::Protobuf::Object.new(hash: "HEAD")))
p commit

