require "gitm"

class LogServer < Gitm::Protobuf::Log::Service
  def get(input, _unused_call)
    r, w = IO.pipe
    repository = input["repository"]
    if repository
      dir_args = ["-C", input["repository"]["path"]]
    end
    system("git", *dir_args, "log", "--format=%H %P", "-n1", input["object"]["hash"], out: w)
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
