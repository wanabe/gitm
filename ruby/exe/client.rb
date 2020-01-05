require "gitm"

stub = Gitm::Protobuf::Log::Stub.new('localhost:50051', :this_channel_is_insecure)
path = ARGV[0]
repository = Gitm::Protobuf::Repository.new(path: path)
iter = Gitm::Protobuf::LogIterator.new(num: 5, repository: repository)
iter = stub.init(iter)

begin
  iter = stub.get(iter)
  p iter.commits.map {|c| c.object["hash"].unpack("H*").first }
end until iter.pointers.empty?
