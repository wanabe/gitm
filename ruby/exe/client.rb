require "gitm"

stub = Gitm::Protobuf::Log::Stub.new('localhost:50051', :this_channel_is_insecure)
iter = stub.init(Gitm::Protobuf::LogIterator.new())
begin
  iter = stub.get(iter)
  p iter.commits.pop
end until iter.pointers.empty?
