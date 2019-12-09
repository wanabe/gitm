require "gitm"

class LogServer < Gitm::Protobuf::Log::Service
  def get(obj, _unused_call)
    r, w = IO.pipe
    system("git", "log", "--format=%H %P", "-n1", obj["hash"], out: w)
    w.close
    hash, parent_hash = r.read.chomp.split(/ +/)
    return Gitm::Protobuf::Commit.new(object: Gitm::Protobuf::Object.new(hash: hash), parent: Gitm::Protobuf::Object.new(hash: parent_hash))
  end
end

s = GRPC::RpcServer.new
s.add_http2_port('0.0.0.0:50051', :this_port_is_insecure)
s.handle(LogServer)
s.run_till_terminated_or_interrupted([1, 'int', 'SIGQUIT'])
