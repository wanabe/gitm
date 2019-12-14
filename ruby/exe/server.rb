require "gitm"

class LogServer < Gitm::Protobuf::Log::Service
  def get(obj, _unused_call)
    r, w = IO.pipe
    system("git", "log", "--format=%H %P", "-n1", obj["hash"], out: w)
    w.close
    hashes = r.read.chomp.split(/ +/)
    object, *parents = hashes.map { |hash| Gitm::Protobuf::Object.new(hash: hash) }
    return Gitm::Protobuf::Commit.new(object: object, parents: parents)
  end
end

s = GRPC::RpcServer.new
s.add_http2_port('0.0.0.0:50051', :this_port_is_insecure)
s.handle(LogServer)
s.run_till_terminated_or_interrupted([1, 'int', 'SIGQUIT'])
